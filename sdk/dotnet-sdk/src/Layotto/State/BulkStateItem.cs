using System;
using System.Collections.Generic;

namespace Layotto.State
{
    public class BulkStateItem
    {
        public string Key { get; set; }

        public Memory<byte> Value { get; set; }

        public string ETag { get; set; }

        public Dictionary<string, string> Metadata { get; set; }

        public string Error { get; set; }
    }
}