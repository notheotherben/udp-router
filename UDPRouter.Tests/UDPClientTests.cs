namespace UDPRouter.Tests
{
    using System;
    using System.Threading.Tasks;
    using UDPRouter.Protocol;
    using Xunit;

    public class UDPClientTests
    {
        [Fact]
        public async Task TestSendAsync()
        {
            var port = new Random().Next(25000, 50000);
            var client = new Client(port);
            var server = new Server(port);


            var data = new byte[] { 1, 2, 3, 4, 5 };
            await client.SendAsync(data);

            var result = await server.ReceiveAsync();
            Assert.NotEmpty(result);
            Assert.Equal(data, result);
        }

        [Fact]
        public async Task TestSendDataAsync()
        {
            var port = new Random().Next(25000, 50000);
            var client = new Client(port);
            var server = new Server(port);

            await client.SendDataAsync(1, 1, "test");

            var packet = await server.ReceivePacketAsync();
            Assert.IsType(typeof(Protocol.DataPacket), packet);
            Assert.Equal("test", ((Protocol.DataPacket)packet).Message);
        }

        [Fact]
        public async Task TestSendControlAsync()
        {
            var port = new Random().Next(25000, 50000);
            var client = new Client(port);
            var server = new Server(port);

            await client.SendControlAsync(1, 1, new Protocol.Route
            {
                Source = 5,
                Dest = 7,
                Port = 10007,
                Cost = 2,
            });

            var packet = await server.ReceivePacketAsync();
            Assert.IsType(typeof(Protocol.ControlPacket), packet);

            var cp = (Protocol.ControlPacket)packet;
            Assert.Equal(5, cp.Route.Source);
            Assert.Equal(7, cp.Route.Dest);
            Assert.Equal(10007, cp.Route.Port);
            Assert.Equal(2, cp.Route.Cost);
        }

    }
}