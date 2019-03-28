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

            var server = new Server(Port);

            while (this.running)
            {
                var packet = await server.ReceivePacketAsync();
                Console.Write($"[{packet.Header.Type}:{packet.Header.Source}->{packet.Header.Dest}] ");

                switch (packet)
                {
                    case ControlPacket p:
                        Console.WriteLine($"{p.Route.Source}->{p.Route.Dest} via :{p.Route.Port} (costs {p.Route.Cost})");
                        break;
                    case DataPacket p:
                        Console.WriteLine(p.Message);
                        break;
                    default:
                        Console.WriteLine("<unrecognized packet type>");
                        break;
                }
            }

            return 0;
        }
    }
}