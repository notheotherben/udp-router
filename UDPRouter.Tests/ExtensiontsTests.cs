namespace UDPRouter.Tests
{
    using UDPRouter.Protocol;
    using Xunit;

    public class ExtensionsTests
    {
        [Fact]
        public void TestConcat()
        {
            var a = new[] { 1, 2, 3, 4 };
            var b = new[] { 5, 6, 7, 8 };

            Assert.Equal(new[] { 1, 2, 3, 4, 5, 6, 7, 8 }, a.Concat(b));
        }
    }
}