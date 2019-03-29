using UDPRouter.Protocol;

namespace UDPRouter.RoutingTable
{
    public class RouteAdapter : Graph.IAdapter<Route>
    {
        public int Cost(Route item)
        {
            return item.Cost;
        }

        public int SourceId(Route item)
        {
            return item.Source;
        }

        public int TargetId(Route item)
        {
            return item.Dest;
        }
    }
}