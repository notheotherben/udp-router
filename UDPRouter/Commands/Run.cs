using System;
using System.Threading.Tasks;
using CommandLine;
using UDPRouter.Protocol;

namespace UDPRouter.Commands
{
    [Verb("run", HelpText = "Run the router process")]
    public class Run : Options
    {
        private bool running = true;

        private void OnCancelKeyPressed(object sender, EventArgs e)
        {
            this.running = false;
        }

        public async Task<int> Execute()
        {
            Console.CancelKeyPress += OnCancelKeyPressed;

            var graph = new Graph.Graph<Route>(ID, new RoutingTable.RouteAdapter());

            await Task.WhenAll(
                RunServer(graph),
                RunRoutePropagator(graph)
            );

            return 0;
        }

        private async Task RunServer(Graph.Graph<Route> graph)
        {
            var routingTable = new RoutingTable.RoutingTable(graph);
            var server = new Server(Port);

            while (this.running)
            {
                var packet = await server.ReceivePacketAsync();

                if (packet.Header.Dest != ID)
                {
                    var route = routingTable[packet.Header.Dest];
                    if (!route.HasValue)
                    {
                        Console.WriteLine($"could not route packet to destination ({packet.Header.Dest}): no known route");
                        continue;
                    }

                    var client = new Client(route.Value.Port);
                    await client.SendAsync(packet.ToBytes());
                    Console.WriteLine($"Forwarded packet from {packet.Header.Source} to {packet.Header.Dest} via port {route.Value.Port}");
                    continue;
                }

                switch (packet)
                {
                    case ControlPacket p:
                        graph.AddOrUpdate(p.Route);
                        routingTable = new RoutingTable.RoutingTable(graph);
                        Console.WriteLine($"Updated routing table");
                        break;
                    case DataPacket p:
                        Console.WriteLine(p.Message);
                        break;
                    default:
                        Console.WriteLine("<unrecognized packet type>");
                        break;
                }
            }
        }

        private async Task RunRoutePropagator(Graph.Graph<Route> graph, int interval = 5000)
        {
            while (this.running)
            {
                await Task.Delay(interval);

                var routingTable = new RoutingTable.RoutingTable(graph);
                foreach (var neighbour in graph.PathsFrom(ID))
                {
                    var client = new Client(neighbour.Port);

                    foreach (var route in routingTable)
                    {
                        await client.SendControlAsync(ID, neighbour.Dest, route);
                    }
                }
            }
        }
    }
}