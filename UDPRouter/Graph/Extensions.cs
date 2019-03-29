using System.Collections.Generic;

namespace UDPRouter.Graph
{
    public static class Extensions
    {
        public static T? Get<K, T>(this Dictionary<K, T> dict, K key)
            where T : struct
        {
            if (dict.TryGetValue(key, out var value))
                return value;
            return null;
        }
    }
}