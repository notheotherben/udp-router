using System;
using CommandLine;

namespace UDPRouter
{
    class Program
    {
        static int Main(string[] args)
        {
            return Parser.Default.ParseArguments<Commands.Run, Commands.Configure, Commands.Send>(args)
                .MapResult(
                    (Commands.Run run) => run.Execute().Result,
                    (Commands.Configure run) => run.Execute().Result,
                    (Commands.Send run) => run.Execute().Result,
                    errs => 1);
        }
    }
}
