using System.Collections.Generic;

namespace Layotto.Configuration
{
    public static class ConfigurationUtil
    {
        public static SubConfigurationResponse ToSubConfigurationResponse(Protocol.SubscribeConfigurationResponse resp)
        {
            var res = new SubConfigurationResponse()
            {
                AppId = resp.AppId,
                StoreName = resp.StoreName,
                Items = new List<ConfigurationItem>()
            };

            foreach (var v in resp.Items)
            {
                var t = new ConfigurationItem()
                {
                    Metadata = new Dictionary<string, string>(v.Metadata),
                    Content = v.Content,
                    Group = v.Group,
                    Key = v.Key,
                    Label = v.Label,
                    Tags = new Dictionary<string, string>(v.Tags)
                };
                res.Items.Add(t);
            }

            return res;
        }
    }
}