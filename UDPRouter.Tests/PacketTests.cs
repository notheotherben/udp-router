using System;
using UDPRouter.Protocol;
using Xunit;

namespace UDPRouter.Tests
{
    public class PacketTests
    {
        [Fact]
        public void TestControlPackets()
        {
            var packet = PacketHeader.ForControl(0, 1);

            Assert.Equal(0, packet.Source);
            Assert.Equal(1, packet.Dest);
            Assert.Equal(16, packet.PayloadSize);

            var bytes = packet.ToBytes();
            Assert.NotNull(bytes);
            Assert.True(bytes.Length > 0);

            var rt = bytes.FromBytes<PacketHeader>();
            Assert.Equal(packet.Type, rt.Type);
            Assert.Equal(packet.Source, rt.Source);
            Assert.Equal(packet.Dest, rt.Dest);
            Assert.Equal(packet.PayloadSize, rt.PayloadSize);
        }

        [Fact]
        public void TestDataPackets()
        {
            var packet = PacketHeader.ForData(0, 1, 17);

            Assert.Equal(0, packet.Source);
            Assert.Equal(1, packet.Dest);
            Assert.Equal(17, packet.PayloadSize);

            var bytes = packet.ToBytes();
            Assert.NotNull(bytes);
            Assert.True(bytes.Length > 0);

            var rt = bytes.FromBytes<PacketHeader>();
            Assert.Equal(packet.Type, rt.Type);
            Assert.Equal(packet.Source, rt.Source);
            Assert.Equal(packet.Dest, rt.Dest);
            Assert.Equal(packet.PayloadSize, rt.PayloadSize);
        }
    }
}
