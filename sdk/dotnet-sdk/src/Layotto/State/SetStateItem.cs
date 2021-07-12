using System;
using System.Collections.Generic;

namespace Layotto.State
{
    public class SetStateItem
    {
        public string Key { get; set; }

        public Memory<byte> Value { get; set; }

        public ETag ETag { get; set; }

        public Dictionary<string, string> Metadata { get; set; }

        public StateOptions Options { get; set; }
    }
}