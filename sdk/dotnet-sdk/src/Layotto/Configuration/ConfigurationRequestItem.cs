using System.Collections.Generic;

namespace Layotto.Configuration
{
    /// <summary>
    /// used for GET,DEL,SUB request
    /// </summary>
    public class ConfigurationRequestItem
    {
        /// <summary>
        /// The name of configuration store.
        /// </summary>
        public string StoreName { get; set; }

        /// <summary>
        /// The application id which
        /// Only used for admin, Ignored and reset for normal client
        /// </summary>
        public string AppId { get; set; }

        /// <summary>
        /// The group of keys.
        /// </summary>
        public string Group { get; set; }

        /// <summary>
        /// The label for keys.
        /// </summary>
        public string Label { get; set; }

        /// <summary>
        /// The keys to get.
        /// </summary>
        public List<string> Keys { get; set; }

        /// <summary>
        /// The metadata which will be sent to configuration store components.
        /// </summary>
        /// <returns></returns>
        public Dictionary<string, string> Metadata { get; set; }
    }
}