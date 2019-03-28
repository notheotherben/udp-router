using System;
using CommandLine;

namespace UDPRouter
{
    public class Options
    {
        [Option("id", Required = true, HelpText = "The unique ID of the router you are connecting to or running.")]
        public int ID { get; set; }

        [Option("port", Required = true, HelpText = "The port on which the router you are connecting to, or running, is listening on.")]
        public int Port { get; set; }
    }
}