using System.Collections.Generic;
using System.Linq;

namespace UDPRouter.Graph
{
    public class Graph<T>
    {
        public Graph(IAdapter<T> adapter)
        {
            Adapter = adapter;
        }

        public IAdapter<T> Adapter { get; set; }

        public List<T> Paths { get; private set; } = new List<T>();

        public IEnumerable<int> Nodes => Paths.Select(Adapter.SourceId).Concat(Paths.Select(Adapter.TargetId)).Distinct();


    }
}