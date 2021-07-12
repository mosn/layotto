using System.Collections.Generic;

namespace Layotto.Configuration
{
    public class SaveConfigurationRequest
    {
        /// <summary>
        /// The name of configuration store.
        /// </summary>
        public string StoreName { get; set; }

        /// <summary>
        /// The application id which
        /// Only used for admin, ignored and reset for normal client
        /// </summary>
        public string AppId { get; set; }

        /// <summary>
        /// The list of configuration items to save.
        /// To delete a exist item, set the key (also label) and let content to be empty
        /// </summary>
        public List<ConfigurationItem> Items { get; set; }

        /// <summary>
        /// The metadata which will be sent to configuration store components.
        /// </summary>
        public Dictionary<string, string> Metadata { get; set; }
    }
}