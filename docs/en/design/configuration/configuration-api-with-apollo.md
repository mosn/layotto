# Implement Configuration API with ctripcorp/apollo
## Goals
Implement [Configuration API(level-2)](https://github.com/dapr/dapr/issues/2988) with ctripcorp/apollo.

## Schema

### From Configuration API to apollo schema

The mapping rule is:

| Params in Configuration API | apollo schema                                                              |
| --------------------------- | -------------------------------------------------------------------------- |
| app_id                      | //ignore it when querying or subscribing,use app_id in config file instead |
| group                       | namespace                                                                  |
| label                       | //append to key. So the actual key stored in apollo will be 'key@$label'   |
| tag                         | // put into another configuration item with json format                    |

The actual key stored in apollo will be `key@$label` and the value will be raw value.

Tags will be stored in a special namespace `sidecar_config_tags`,

with key=`group@$key@$label` and value=

```json
{
  "tag1": "tag1value",
  "tag2": "tag2value"
}
```


**Q: Why not store value and tags in a single configuration item to reduce times of queries,like:**

```json
{
    "value":"raw value",
    "tags":{
        "tag1":"tag1value",
        "tag2":"tag2value"
    }
}
```

A: Legacy systems using apollo can't migrate to our sidecar if we design like this.

### How to migrate legacy systems

1. Get and subscribe APIs are compatible.Users can easily put legacy systems onto our sidecar if they don't use save/delete APIs.Just keep `label` field blank in config.json,and the sidecar will use the raw key instead of `key@$label` to interact with apollo.

2. Save/delete APIs might be incompatible.The sidecar use fixed `cluster` field configurated in config.json and fixed `env` field in code,which means users can't pass `cluster` and `env` field as a parameter for save/delete API when they want to change some configuration items with other appid.

### config.json for sidecar

```json
{
  "config_store": {
    "type": "apollo",
    "address": [
      "http://106.54.227.205:8080"
    ],
    "metadata": {
      "app_id": "testApplication_yang",
      "cluster": "dev",
      "namespace_name": "dubbo,product.joe",
      "is_backup_config": true,
      "secret": "6ce3ff7e96a24335a9634fe9abca6d51"
    }
  }
}
```


## API

### Which Go SDK for apollo should I Use?

We are using the official maintained [sdk](https://github.com/apolloconfig/agollo), the others sdks you can find in this [repo](https://www.apolloconfig.com/#/zh/usage/third-party-sdks-user-guide).

Some problems with the sdk:
1. Users must declare all namespaces in AppConfig before connecting to the server and constructing a client,like:

```go
	c := &config.AppConfig{
		AppID:          "testApplication_yang",
		Cluster:        "dev",
		IP:             "http://106.54.227.205:8080",
		NamespaceName:  "dubbo",  // declare before connecting to the server
		IsBackupConfig: true,
		Secret:         "6ce3ff7e96a24335a9634fe9abca6d51",
	}
	client,err:=agollo.StartWithConfig(func() (*config.AppConfig, error) {
		return c, nil
	})
```

2. Nowhere to configurate `env` field.

3. No save/delete API.

4. No bulk query API.

5. No sure about the concurrent safety.

6. Nowhere to configurate or use [Apollo Meta Server](https://www.apolloconfig.com/#/zh/usage/java-sdk-user-guide?id=_122-apollo-meta-server)

7. Not sure about the consistency property between local cache and backend database.

8. The two operations(set kv+set tags) are not transaction,which may cause inconsistency.

### Which apollo sdk API should I use?

| Configuration API      | apollo sdk API                                                                                                                                                                                                                                                                                                 |
| ---------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| GetConfiguration       | cache := client.GetConfigCache(c.NamespaceName)  <br> value,_ := client.Get("key")                                                                                                                                                                                                                             |
| SaveConfiguration      | use open API via http. [update](https://www.apolloconfig.com/#/zh/usage/apollo-open-api-platform?id=_3211-%e4%bf%ae%e6%94%b9%e9%85%8d%e7%bd%ae%e6%8e%a5%e5%8f%a3) + [commit](https://www.apolloconfig.com/#/zh/usage/apollo-open-api-platform?id=_3213-%e5%8f%91%e5%b8%83%e9%85%8d%e7%bd%ae%e6%8e%a5%e5%8f%a3) |
| DeleteConfiguration    | use open API via http. [delete](https://www.apolloconfig.com/#/zh/usage/apollo-open-api-platform?id=_3212-%e5%88%a0%e9%99%a4%e9%85%8d%e7%bd%ae%e6%8e%a5%e5%8f%a3) + [commit](https://www.apolloconfig.com/#/zh/usage/apollo-open-api-platform?id=_3213-%e5%8f%91%e5%b8%83%e9%85%8d%e7%bd%ae%e6%8e%a5%e5%8f%a3) |
| SubscribeConfiguration | https://github.com/apolloconfig/agollo/wiki/%E7%9B%91%E5%90%AC%E5%8F%98%E6%9B%B4%E4%BA%8B%E4%BB%B6                                                                                                                                                                                                             |

### How to subscribe all config changes of the specified app

Subscribe all namespaces declared in config.json