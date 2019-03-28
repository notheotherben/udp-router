namespace UDPRouter.Protocol
{
    public enum PacketType : byte
    {
        Control = 0,
        Data = 1
    }

    public class Packet
    {
        public Packet(PacketHeader header) => Header = header;

        public PacketHeader Header { get; private set; }

    }

    public class ControlPacket : Packet
    {
        public ControlPacket(PacketHeader header, Route route)
            : base(header) => Route = route;

        public Route Route { get; private set; }
    }

    public class DataPacket : Packet
    {
        public DataPacket(PacketHeader header, string message)
            : base(header) => Message = message;

        public string Message { get; private set; }
    }
}