using System.Net.Sockets;
using System.Threading.Tasks;

namespace UDPRouter.Protocol
{
    public class Client
    {
        private readonly UdpClient client;

        public Client(int port) => this.client = new UdpClient("127.0.0.1", port);

        ~Client() => client.Close();

        public async Task SendAsync(byte[] data) => await this.client.SendAsync(data, data.Length);

        public async Task SendAsync(PacketHeader header, byte[] payload) => await this.SendAsync(header.ToBytes().Concat(payload));

        public async Task SendControlAsync(int source, int dest, Route route) => await this.SendAsync(PacketHeader.ForControl(source, dest), route.ToBytes());

        public async Task SendDataAsync(int source, int dest, byte[] data) => await this.SendAsync(PacketHeader.ForData(source, dest, data.Length), data);

        public async Task SendDataAsync(int source, int dest, string data) => await this.SendDataAsync(source, dest, System.Text.Encoding.UTF8.GetBytes(data));
    }
}