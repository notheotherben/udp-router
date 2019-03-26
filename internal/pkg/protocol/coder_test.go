package protocol_test

import (
	"bytes"
	"testing"

	"github.com/spartan563/udp-router/internal/pkg/protocol"

	"github.com/stretchr/testify/assert"
)

func TestCoder(t *testing.T) {
	tcs := []protocol.Packet{
		protocol.Packet{
			PacketHeader: protocol.PacketHeader{
				Type:    protocol.ControlPacketType,
				Subtype: protocol.PathAdvertisementSubtype,
				Source:  protocol.Address(1),
				Dest:    protocol.Address(1),
			},
			Payload: protocol.PathAdvertisement{
				Source: protocol.Address(1),
				Dest:   protocol.Address(5),
				Port:   10005,
				Cost:   5,
			},
		},
		protocol.Packet{
			PacketHeader: protocol.PacketHeader{
				Type:    protocol.DataPacketType,
				Source:  protocol.Address(1),
				Dest:    protocol.Address(1),
			},
			Payload: protocol.DataPayload{
				Length: 4,
				Data: []byte{1,2,3,4},
			},
		},
	}

	for _, tc := range tcs {
		buf := bytes.NewBuffer([]byte{})
		assert.NotNil(t, buf)

		enc := protocol.NewEncoder(buf)
		assert.NotNil(t, enc)
		assert.Nil(t, enc.Encode(&tc))

		dec := protocol.NewDecoder(buf)
		assert.NotNil(t, dec)

		var p protocol.Packet
		assert.Nil(t, dec.Decode(&p))

		assert.Equal(t, tc, p)
	}
}
