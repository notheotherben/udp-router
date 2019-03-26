package router

import (
	"bytes"
	"fmt"
	"net"

	"github.com/sirupsen/logrus"

	"github.com/spartan563/udp-router/internal/pkg/protocol"

	"github.com/spartan563/udp-router/internal/pkg/netmap"
)

type Router struct {
	address protocol.Address
	netmap  *netmap.Map
	conn    net.PacketConn
}

func New(listen string, address protocol.Address) (*Router, error) {
	r := &Router{
		address: address,
		netmap:  netmap.New(address),
	}

	conn, err := net.ListenPacket("udp", listen)
	if err != nil {
		return nil, err
	}

	r.conn = conn
	return r, nil
}

func (r *Router) Run() {
	r.receiveLoop()
}

func (r *Router) send(packet *protocol.Packet, via netmap.Route) error {
	buf := bytes.NewBuffer([]byte{})
	if err := protocol.NewEncoder(buf).Encode(packet); err != nil {
		return err
	}

	dst, err := net.ResolveUDPAddr("udp", fmt.Sprintf("127.0.0.1:%d", via.Port))
	if err != nil {
		return err
	}

	_, err = r.conn.WriteTo(buf.Bytes(), dst)
	if err != nil {
		return err
	}

	return nil
}

func (r *Router) receiveLoop() {
	buf := make([]byte, 1024)
	for {
		n, _, err := r.conn.ReadFrom(buf)

		if err != nil {
			logrus.WithError(err).Error("failed to read packet")
			return
		}

		//from := addr.(*net.UDPAddr).Port
		// TODO: We need to propagate the source of packets to
		//       the onPacket function

		if n > 0 {
			go r.onPacket(buf[:n])
		}
	}
}

func (r *Router) onPacket(data []byte) {
	buf := bytes.NewBuffer(data)
	var packet protocol.Packet
	if err := protocol.NewDecoder(buf).Decode(&packet); err != nil {
		logrus.WithError(err).Warning("failed to decode packet")
		return
	}

	if packet.Dest != r.address {
		// Forward packet
		viaRoute, ok := r.netmap.Route(packet.Dest)
		if !ok {
			logrus.WithFields(logrus.Fields{
				"source":  packet.Source,
				"dest":    packet.Dest,
				"type":    packet.Type,
				"subtype": packet.Subtype,
			}).Warning("unable to route packet to destination")
		}

		if err := r.send(&packet, viaRoute.Route()); err != nil {
			logrus.WithError(err).Warning("failed to send packet to destination")
		}

		return
	}

	switch packet.Type {
	case protocol.ControlPacketType:
		r.onControlPacket(&packet)
	case protocol.DataPacketType:
		r.onDataPacket(&packet)
	}
}

func (r *Router) onDataPacket(packet *protocol.Packet) {
	logrus.WithFields(logrus.Fields{
		"source":  packet.Source,
		"payload": string(packet.Payload.(protocol.DataPayload).Data),
	}).Info("received message")
}

func (r *Router) onControlPacket(packet *protocol.Packet) {
	switch packet.Subtype {
	case protocol.PathAdvertisementSubtype:
		update := packet.Payload.(protocol.PathAdvertisement)

		logrus.WithFields(logrus.Fields{
			"source":  packet.Source,
			"payload": fmt.Sprintf("%d -> %d costs %d (port %d)", update.Source, update.Dest, update.Cost, update.Port),
		}).Info("received network update")

		diff := r.netmap.Update(netmap.Route{
			Source: update.Source,
			Dest:   update.Dest,
			Port:   update.Port,
			Cost:   update.Cost,
		})

		ns := r.netmap.Neighbours()
		for _, vr := range diff {
			for _, n := range ns {
				p := protocol.Packet{
					PacketHeader: protocol.PacketHeader{
						Source:  r.address,
						Dest:    n.Dest,
						Type:    protocol.ControlPacketType,
						Subtype: protocol.PathAdvertisementSubtype,
					},
					Payload: &protocol.PathAdvertisement{
						Source: vr.Source,
						Dest:   vr.Dest,
						Cost:   vr.EstCost,
						Port:   vr.Via.Port,
					},
				}

				if err := r.send(&p, n); err != nil {
					logrus.WithFields(logrus.Fields{
						"neighbour": n.Dest,
						"port":      n.Port,
					}).WithError(err).Warning("failed to propagate route update")
				}
			}
		}
	}
}
