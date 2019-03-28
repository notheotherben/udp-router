using System.Runtime.InteropServices;

namespace UDPRouter.Protocol
{
    [StructLayout(LayoutKind.Sequential)]
    public struct PacketHeader
    {
        public PacketType Type;

        public int Source;

        public int Dest;

        public int PayloadSize;

        public static PacketHeader ForControl(int source, int dest)
        {
            return new PacketHeader
            {
                Type = PacketType.Control,
                Source = source,
                Dest = dest,
                PayloadSize = Marshal.SizeOf(typeof(Route)),
            };
        }

        public static PacketHeader ForData(int source, int dest, int dataSize)
        {
            return new PacketHeader
            {
                Type = PacketType.Data,
                Source = source,
                Dest = dest,
                PayloadSize = dataSize,
            };
        }
    }
}