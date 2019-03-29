using System.Collections.Generic;
using System.Linq;
using Xunit;

namespace UDPRouter.Tests
{
    public class GraphTests
    {
        [Fact]
        public void TestBasicGraph()
        {
            var graph = new Graph.Graph<TestPath>(0, new TestAdapter());

            graph.AddOrUpdate(new TestPath
            {
                From = 0,
                To = 1,
                Cost = 5,
            });

            graph.AddOrUpdate(new TestPath
            {
                From = 1,
                To = 2,
                Cost = 3,
            });

            graph.AddOrUpdate(new TestPath
            {
                From = 0,
                To = 2,
                Cost = 12,
            });

            Assert.Equal(2, graph.Path(2).Count());
            Assert.Equal(8, graph.Path(2).Select(c => c.Cost).Sum());
            Assert.Equal(new[] { 5, 3 }, graph.Path(2).Select(c => c.Cost));
            Assert.Equal(new[] { 0, 1 }, graph.Path(2).Select(c => c.From));
            Assert.Equal(new[] { 1, 2 }, graph.Path(2).Select(c => c.To));

            graph.AddOrUpdate(new TestPath
            {
                From = 0,
                To = 2,
                Cost = 1,
            });

            Assert.Equal(1, graph.Path(2).Count());
            Assert.Equal(1, graph.Path(2).Select(c => c.Cost).Sum());
            Assert.Equal(new[] { 1 }, graph.Path(2).Select(c => c.Cost));
            Assert.Equal(new[] { 0 }, graph.Path(2).Select(c => c.From));
            Assert.Equal(new[] { 2 }, graph.Path(2).Select(c => c.To));
        }

        private class TestPath
        {
            public int From { get; set; }
            public int To { get; set; }
            public int Cost { get; set; }
        }

        private class TestAdapter : Graph.IAdapter<TestPath>
        {
            public int Cost(TestPath item)
            {
                return item.Cost;
            }

            public int SourceId(TestPath item)
            {
                return item.From;
            }

            public int TargetId(TestPath item)
            {
                return item.To;
            }
        }
    }
}