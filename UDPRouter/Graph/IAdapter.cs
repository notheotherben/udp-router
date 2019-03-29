using System.Collections.Generic;

namespace UDPRouter.Graph
{
    public interface IAdapter<T>
    {
        int SourceId(T item);

        int TargetId(T item);

        int Cost(T item);
    }
}