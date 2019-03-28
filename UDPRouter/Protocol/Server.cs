using System;
using System.Diagnostics;
using System.Net.Sockets;
using System.Runtime.InteropServices;
using System.Threading.Tasks;

namespace UDPRouter.Protocol
{
    public class Server
    {
        private readonly UdpClient client;

        public Server(int port)
        {
            this.client = new UdpClient(port);
        }

        ~Server() => client.Close();

        public async Task<byte[]> ReceiveAsync()
        {
            var msg = await client.ReceiveAsync();
            return msg.Buffer;
        }

        public async Task<Packet> ReceivePacketAsync()
        {
            var result = await this.ReceiveAsync();

            var header = result.FromBytes<PacketHeader>();
            var headerSize = Marshal.SizeOf(header);

            Debug.Assert(headerSize + header.PayloadSize == result.Length);

            var payload = result.AsSpan(headerSize, header.PayloadSize).ToArray();

            switch (header.Type)
            {
                case PacketType.Control:
                    return new ControlPacket(header, payload.FromBytes<Route>());
                case PacketType.Data:
                    return new DataPacket(header, System.Text.Encoding.UTF8.GetString(payload));
                default:
                    throw new NotImplementedException("unrecognized packet type");
            }
        }
    }
}