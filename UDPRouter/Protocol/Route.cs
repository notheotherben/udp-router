using System.Runtime.InteropServices;

namespace UDPRouter.Protocol
{
    [StructLayout(LayoutKind.Sequential)]
    public struct Route
    {
        public int Source;

        public int Dest;

        public int Port;

        public int Cost;
    }
}