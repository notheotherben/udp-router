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

	"github.com/sirupsen/logrus"
	"github.com/spartan563/udp-router/internal/pkg/protocol"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// sendCmd represents the send command
var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "Send a message over the network",
	Long: `Injects a message into the network through a specific router, allowing it to be
	routed to its destination.
	
	You will be required to provide the destination and payload, with the source defaulting to
	the ID of the router that it is injected at.`,
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

		payload, err := cmd.Flags().GetString("payload")
		if err != nil {
			logrus.WithError(err).Fatal()
		}

		log := logrus.WithFields(logrus.Fields{
			"id":      id,
			"port":    port,
			"src":     src,
			"dest":    dest,
			"payload": payload,
		})

		log.Info("Sending message via router")

		packet := &protocol.Packet{
			PacketHeader: protocol.PacketHeader{
				Type:   protocol.DataPacketType,
				Source: protocol.Address(src),
				Dest:   protocol.Address(dest),
			},
			Payload: &protocol.DataPayload{
				Data:   []byte(payload),
				Length: len(payload),
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
			logrus.WithError(err).Error("failed to send message to router")
			return
		}

		logrus.Info("sent message to router")
	},
}

func init() {
	rootCmd.AddCommand(sendCmd)

	sendCmd.Flags().Int("src", 0, "id of source router")

	sendCmd.Flags().Int("dest", 0, "id of destination router")
	sendCmd.MarkFlagRequired("dest")

	sendCmd.Flags().String("payload", "", "message to transmit")
	sendCmd.MarkFlagRequired("payload")
}
