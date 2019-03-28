using System;
using System.Threading.Tasks;
using CommandLine;
using UDPRouter.Protocol;

namespace UDPRouter.Commands
{
    [Verb("configure", HelpText = "Configure a router on the network with an updated route entry")]
    public class Configure : Options
    {
        [Option("src", Required = true, HelpText = "The unique ID of the router from which the route originages.")]
        public int Source { get; set; }

        [Option("dest", Required = true, HelpText = "The unique ID of the router at which the route terminates.")]
        public int Dest { get; set; }

        [Option("cost", Required = true, HelpText = "The estimated route traversal cost.")]
        public int Cost { get; set; }

        [Option("dport", Required = true, HelpText = "The port to which the packet should be sent for this route.")]
        public int DestPort { get; set; }

        public async Task<int> Execute()
        {
            var client = new Client(this.Port);
            await client.SendControlAsync(this.ID, this.ID, new Route
            {
                Source = Source,
                Dest = Dest,
                Port = Port,
                Cost = Cost,
            });

            return 0;
        }
    }
}