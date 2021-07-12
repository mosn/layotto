using System;

namespace Layotto.Configuration
{
    public class WatchResponse
    {
        public SubConfigurationResponse Item { get; set; }
        public Exception Error { get; set; }
    }
}