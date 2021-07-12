using System;
using System.Collections.Generic;

namespace Layotto.State
{
    public class StateItem
    {
        public string Key { get; set; }

        public Memory<byte> Value { get; set; }

        public string ETag { get; set; }

        public Dictionary<string, string> Metadata { get; set; }
    }
}