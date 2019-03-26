// Copyright © 2019 Benjamin Pannell <benjamin@pannell.dev>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"bytes"
	"fmt"
	"net"

	"github.com/spartan563/udp-router/internal/pkg/protocol"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// configureCmd represents the configure command
var configureRouteCmd = &cobra.Command{
	Use:   "route",
	Short: "Configure a route",
	Long: `Configures a router by injecting a control packet on its interface, allowing
	you to advertise new routes or update the path cost of existing ones.`,
	Run: func(cmd *cobra.Command, args []string) {
		id := viper.GetInt("id")
		port := viper.GetInt("port")

		src, err := cmd.Flags().GetInt("src")
		if err != nil {
			src = id
		}

		dest, err := cmd.Flags().GetInt("dest")
		if err != nil {
			logrus.WithError(err).Fatal()
		}

		destPort, err := cmd.Flags().GetInt("dport")
		if err != nil {
			logrus.WithError(err).Fatal()
		}

		cost, err := cmd.Flags().GetInt("cost")
		if err != nil {
			logrus.WithError(err).Fatal()
		}

		log := logrus.WithFields(logrus.Fields{
			"id":    id,
			"port":  port,
			"src":   src,
			"dest":  dest,
			"dport": destPort,
			"cost":  cost,
		})

		log.Info("Configuring route advertisement")

		packet := &protocol.Packet{
			PacketHeader: protocol.PacketHeader{
				Type:    protocol.ControlPacketType,
				Subtype: protocol.PathAdvertisementSubtype,
				Source:  protocol.Address(id),
				Dest:    protocol.Address(id),
			},
			Payload: &protocol.PathAdvertisement{
				Source: protocol.Address(src),
				Dest:   protocol.Address(dest),
				Port:   destPort,
				Cost:   cost,
			},
		}

		buf := bytes.NewBuffer([]byte{})
		if err := protocol.NewEncoder(buf).Encode(packet); err != nil {
			logrus.WithError(err).Error("failed to serialize packet")
			return
		}

		conn, err := net.Dial("udp", fmt.Sprintf("127.0.0.1:%d", port))
		if err != nil {
			logrus.WithError(err).Error("failed to connect to router")
			return
		}

		_, err = conn.Write(buf.Bytes())
		if err != nil {
			logrus.WithError(err).Error("failed to send command to router")
			return
		}

		logrus.Info("sent command to router")
	},
}

func init() {
	configureCmd.AddCommand(configureRouteCmd)

	configureRouteCmd.Flags().Int("src", 0, "id of source router")

	configureRouteCmd.Flags().Int("dest", 0, "id of destination router")
	configureRouteCmd.MarkFlagRequired("dest")

	configureRouteCmd.Flags().Int("dport", 0, "port of destination router")
	configureRouteCmd.MarkFlagRequired("dport")

	configureRouteCmd.Flags().Int("cost", 0, "cost of path traversal")
	configureRouteCmd.MarkFlagRequired("cost")
}
