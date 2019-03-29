namespace UDPRouter.RoutingTable
{
    using System.Collections;
    using System.Collections.Generic;
    using System.Linq;
    using UDPRouter.Graph;
    using UDPRouter.Protocol;

    public class RoutingTable : IEnumerable<Route>
    {
        public RoutingTable(Graph.Graph<Route> graph)
        {
            foreach (var node in graph.Nodes)
            {
                if (node == graph.Root) continue;

                var path = graph.Path(node).ToArray();
                if (!path.Any()) continue;

                routes[node] = new Route
                {
                    Source = graph.Root,
                    Dest = node,
                    Cost = path.Select(p => p.Cost).Sum(),
                    Port = path.First().Port,
                };
            }
        }

        private Dictionary<int, Route> routes = new Dictionary<int, Route>();

        public Route? this[int dest] => routes.Get(dest);

        public IEnumerator<Route> GetEnumerator() => routes.Values.GetEnumerator();

        IEnumerator IEnumerable.GetEnumerator() => routes.Values.GetEnumerator();
    }
}