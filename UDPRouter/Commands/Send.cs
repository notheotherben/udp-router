using System;
using System.Threading.Tasks;
using CommandLine;
using UDPRouter.Protocol;

namespace UDPRouter.Commands
{
    [Verb("send", HelpText = "Send some data over the network")]
    public class Send : Options
    {
        [Option("src", Required = true, HelpText = "The unique ID of the router at which this packet originated.")]
        public int Source { get; set; }

        [Option("dest", Required = true, HelpText = "The unique ID of the router to which this packet should be delivered.")]
        public int Dest { get; set; }


        [Option("msg", Required = true, HelpText = "The message which should be transmitted to the destination router.")]
        public string Message { get; set; }

        public async Task<int> Execute()
        {
            if (this.Message.Length >= 100)
                this.Message = this.Message.Substring(0, 99);

            var client = new Client(this.Port);
            await client.SendDataAsync(this.Source, this.Dest, this.Message);

            return 0;
        }
    }
}