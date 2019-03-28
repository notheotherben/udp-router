using System.IO;
using System.Runtime.InteropServices;
using System.Threading.Tasks;

namespace UDPRouter.Protocol
{
    public static class Extensions
    {
        public static byte[] ToBytes<T>(this T structure) where T : struct
        {
            var size = Marshal.SizeOf(structure);
            var memPtr = Marshal.AllocCoTaskMem(size);

            try
            {
                Marshal.StructureToPtr(structure, memPtr, false);

                var bytes = new byte[size];
                Marshal.Copy(memPtr, bytes, 0, size);
                return bytes;
            }
            finally
            {
                Marshal.FreeCoTaskMem(memPtr);
            }
        }

        public static T[] Concat<T>(this T[] a1, T[] a2)
        {
            var a3 = new T[a1.Length + a2.Length];
            for (var i = 0; i < a1.Length; i++)
                a3[i] = a1[i];

            for (var i = 0; i < a2.Length; i++)
                a3[a1.Length + i] = a2[i];

            return a3;
        }

        public static T FromBytes<T>(this byte[] data) where T : struct
        {
            var handle = GCHandle.Alloc(data, GCHandleType.Pinned);

            try
            {
                return (T)Marshal.PtrToStructure(handle.AddrOfPinnedObject(), typeof(T));
            }
            finally
            {
                handle.Free();
            }
        }

        public static async Task<T> FromStreamAsync<T>(this Stream stream) where T : struct
        {
            var size = Marshal.SizeOf(typeof(T));
            var data = new byte[size];
            var offset = 0;

            while (offset < size)
            {
                offset += await stream.ReadAsync(data, offset, size - offset);
            }

            return data.FromBytes<T>();
        }
    }
}