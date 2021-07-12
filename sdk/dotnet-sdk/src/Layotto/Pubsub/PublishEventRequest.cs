using System;
using System.Collections.Generic;

namespace Layotto.Pubsub
{
    public class PublishEventRequest
    {
        public string PubsubName { get; set; }

        public string Topic { get; set; }

        public Memory<byte> Data { get; set; }

        public string DataContentType { get; set; }

        public Dictionary<string, string> Metadata { get; set; }
    }
}