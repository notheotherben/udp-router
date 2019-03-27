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
	"fmt"
	"os"
	"os/signal"

	"github.com/spartan563/udp-router/internal/pkg/netmap"

	"github.com/spartan563/udp-router/internal/pkg/protocol"
	"github.com/spartan563/udp-router/internal/pkg/router"
	"github.com/spf13/viper"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Starts a routing daemon",
	Long: `Begins running a routing daemon on the provided listening address and with the
	provided unique identifier.
	
	This router may then be configured using the configure command to advertise routes and
	join the network.`,
	Run: func(cmd *cobra.Command, args []string) {
		id := viper.GetInt("id")
		port := viper.GetInt("port")

		log := logrus.WithFields(logrus.Fields{
			"id":   id,
			"port": port,
		})

		log.Info("Starting router")

		rtr, err := router.New(fmt.Sprintf(":%d", port), protocol.Address(id))
		if err != nil {
			logrus.WithError(err).Error("failed to start router")
			return
		}

		routes := []netmap.Route{}
		if err := viper.UnmarshalKey("routes", &routes); err != nil {
			logrus.WithError(err).Error("failed to configure default routes")
			return
		}

		for _, route := range routes {
			rtr.AddRoute(route)
		}

		rtr.Run()

		sigCh := make(chan os.Signal)
		signal.Notify(sigCh, os.Kill)

		<-sigCh

		rtr.Shutdown()
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
