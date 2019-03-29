using System;
using System.Collections.Generic;
using System.Linq;

namespace UDPRouter.Graph
{
    public class Graph<T>
    {
        public Graph(int root, IAdapter<T> adapter)
        {
            Adapter = adapter;
            Root = root;
        }

        public int Root { get; private set; }

        public IAdapter<T> Adapter { get; private set; }

        public IEnumerable<int> Nodes => Paths.Select(Adapter.SourceId).Concat(Paths.Select(Adapter.TargetId)).Distinct();

        private List<T> Paths { get; set; } = new List<T>();

        private Dictionary<int, int> ShortestPaths { get; set; } = new Dictionary<int, int>();


        public void AddOrUpdate(T path)
        {
            foreach (var item in Paths)
            {
                if (Adapter.SourceId(item) != Adapter.SourceId(path))
                    continue;

                if (Adapter.TargetId(item) != Adapter.TargetId(path))
                    continue;

                if (Adapter.Cost(item) < Adapter.Cost(path))
                    return;

                Paths.Remove(item);
                break;
            }

            Paths.Add(path);
            Update();
        }

        public IEnumerable<T> PathsFrom(int source) => Paths.Where(p => Adapter.SourceId(p) == source);

        public IEnumerable<T> PathsTo(int target) => Paths.Where(p => Adapter.TargetId(p) == target);

        private void Update()
        {
            // Shortest path to ourselves is obviously 0
            ShortestPaths[Root] = 0;

            foreach (var node in Nodes)
            {
                foreach (var path in Paths.Where(p => Adapter.SourceId(p) == node))
                {
                    var estCost = (ShortestPaths.Get(Adapter.SourceId(path)) ?? int.MaxValue);
                    if (estCost < int.MaxValue)
                        estCost += Adapter.Cost(path);

                    if (estCost < (ShortestPaths.Get(Adapter.TargetId(path)) ?? int.MaxValue))
                    {
                        ShortestPaths[Adapter.TargetId(path)] = estCost;
                    }
                }
            }

            foreach (var path in Paths)
            {
                if ((ShortestPaths.Get(Adapter.SourceId(path)) ?? int.MaxValue) + Adapter.Cost(path) < (ShortestPaths.Get(Adapter.TargetId(path)) ?? int.MaxValue))
                    throw new InvalidOperationException("graph contains negative weight cycle");
            }
        }

        public IEnumerable<T> Path(int to)
        {
            if (to == Root) yield break;

            var viaPaths = Paths.Where(p =>
            {
                if (Adapter.TargetId(p) != to) return false;

                if ((ShortestPaths.Get(Adapter.SourceId(p)) ?? int.MaxValue) + Adapter.Cost(p) != (ShortestPaths.Get(to) ?? int.MaxValue)) return false;
                return true;
            });

            if (!viaPaths.Any())
                throw new InvalidOperationException("unable to find route to target");


            var viaPath = viaPaths.First();
            foreach (var path in Path(Adapter.SourceId(viaPath)))
            {
                yield return path;
            }

            yield return viaPath;
        }
    }
}