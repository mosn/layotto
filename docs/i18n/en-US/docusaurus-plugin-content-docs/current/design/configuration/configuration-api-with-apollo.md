# Implementing Configuration API with ctripcorp/apollo

## Objective

Implementation of [Configuration API(level-2)](https://github.com/dapr/dapr/issues/2988) with ctripcorp/apollo

## Schema

### Configuration API map apollo schema

The mapping rules are as follows:

| Configuration API Parameters | apollo schema                                                                                                                  |
| ---------------------------- | ------------------------------------------------------------------------------------------------------------------------------ |
| app_id  | ///ignore this parameter when searching for or subscribed to it, use app_id in configuration file instead |
| Group                        | name                                                                                                                           |
| label                        | //追加到key的后面. 在apollo中key实际存储的格式为'key@$label'                                                      |
| tag                          | // Put an item in the configuration in json                                                                                    |

在apollo中key实际存储的格式为`key@$label`，value为原始值。

Tags将以key=`group@$key@$label` and value=

```json
LO
  "tag1": "tag1value",
  "tag2": "tag2value"
}
```

Format is stored in a specific namespace `sidecar_config_tags`

**Q: Why don't you have an item to store values and tags to reduce the number of queries, e.g. :**

```json
{
    "value":"raw value",
    "tags":{
        "tag1":"tag1value",
        "tag2":"tag2value"
    }
}
```

A: If this design is adopted, the historical legacy of the use of apollo will not be able to migrate to sidecar

### How to migrate historical legacy projects

1. Get/subscripbe API is compatible.Without save/delete API, users can easily migrate legacy projects to sidecar.只需要在config.json文件中保持`label`字段为空，sidecar将会使用原始key来替代`key@$label`与apollo进行交互。
2. Save/delete API may not be compatible.sidecar uses fixed `cluster` and `env` fields in config.json.This means that users cannot use `cluster` and `env` fields as parameters for save/deleteAPI if they want to change some configuration items with other appids.

### sidecar config.json file

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

### Which Apollo SDK should I use?

Currently using official maintenance[SDK](https://github.com/apolloconfig/agollo), other Apollo's third party SDK can be found from[链接](https://www.apolloconfig.com/#/usage/third-party-sdks-user-guide).

Use some of this sdk's attention:

1. Before connecting to the server and creating the client, the user must declare all namespace in the AppConfig, as shown below in：

```go
	c := &config. ppConfigL
		AppID: "testApplication_yang",
		Cluster: "dev",
		IP: "http://106. 4.227. 05:8080",
		NamespaceName: "dubbo", // Statement
		IsBackupConfig: true,
		Secret: "6ce3ff7e96a24335a9634fe9abca6d51",
	}
	client, err := agolloConfig. tartWithConfig(func() (*config.AppConfig, error) {
		return c, nil
	})
```

1. You cannot configure the `env` field.

2. No Save/delete API.

3. No bulk queries for API.

4. Security cannot be ensured in parallel.

5. Configure or don't use [Apollo Meta Server](https://www.apollocconfig.com/#/usage/java-sdk-user-guide?id=_122-apollo-meta-server).

6. There is a problem setting kv and tags. There is no transaction.

### How I should use the Apollo SDK API

| Configuration API       | apollo sdk API                                                                                                                                                                                                                                                                                                             |
| ----------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| GetConfiguration        | cache := client.GetConfigCache(c.NamespaceName)value,_:= client.Get("key")                                                                                                      |
| SaveConfiguration       | 通过http. [update](https://www.apolloconfig.com/#/zh/usage/apollo-open-api-platform?id=_3211-%e4%bf%ae%e6%94%b9%e9%85%8d%e7%bd%ae%e6%8e%a5%e5%8f%a3) + [commit](https://www.apolloconfig.com/#/zh/usage/apollo-open-api-platform?id=_3213-%e5%8f%91%e5%b8%83%e9%85%8d%e7%bd%ae%e6%8e%a5%e5%8f%a3) 使用open API |
| DeleteConfiguration     | 使用http. [delete](https://www.apolloconfig.com/#/zh/usage/apollo-open-api-platform?id=_3212-%e5%88%a0%e9%99%a4%e9%85%8d%e7%bd%ae%e6%8e%a5%e5%8f%a3) + [commit](https://www.apolloconfig.com/#/zh/usage/apollo-open-api-platform?id=_3213-%e5%8f%91%e5%b8%83%e9%85%8d%e7%bd%ae%e6%8e%a5%e5%8f%a3) 使用open API |
| SubscripbeConfiguration | https://github.com/apolloconfig/agollo/wiki/%E7%9B%91%E5%90%AC%E5%8F%98%E6%9B%B4%E4%BA%8B%E4%BB%B6                                                                                                                                                                                         |

### How to subscribe to all configuration changes for the specified app

Subscribe to all declared namespaces in config.json file.
