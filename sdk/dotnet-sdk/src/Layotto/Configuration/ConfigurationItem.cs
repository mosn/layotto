using System.Collections.Generic;

namespace Layotto.Configuration
{
    public class ConfigurationItem
    {
        /// <summary>
        /// Required. The key of configuration item
        /// </summary>
        public string Key { get; set; }

        /// <summary>
        /// The content of configuration item
        /// Empty if the configuration is not set, including the case that the configuration is changed from value-set to
        /// value-not-set.
        /// </summary>
        public string Content { get; set; }

        /// <summary>
        /// The group of configuration item.
        /// </summary>
        public string Group { get; set; }

        /// <summary>
        /// The label of configuration item.
        /// </summary>
        public string Label { get; set; }

        /// <summary>
        /// The tag list of configuration item.
        /// </summary>
        public Dictionary<string, string> Tags { get; set; }

        /// <summary>
        /// The metadata which will be passed to configuration store component.
        /// </summary>
        /// <returns></returns>
        public Dictionary<string, string> Metadata { get; set; }
    }
}