# 使用ctripcorp/apollo实现Configuration API
## 目标
使用ctripcorp/apollo实现[Configuration API(level-2)](https://github.com/dapr/dapr/issues/2988)

## Schema

### Configuration API映射apollo schema

映射规则如下:

| Configuration API参数 | apollo schema                                                |
| --------------------- | ------------------------------------------------------------ |
| app_id                | //当查询或者订阅时忽略该参数，使用配置文件中的app_id替代     |
| group                 | namespace                                                    |
| label                 | //追加到key的后面. 在apollo中key实际存储的格式为'key@$label' |
| tag                   | //以json格式置于配置中的某一项                               |

在apollo中key实际存储的格式为`key@$label`，value为原始值。

Tags将以key=`group@$key@$label` and value=

```json
{
  "tag1": "tag1value",
  "tag2": "tag2value"
}
```

的格式存储在特定的namespace `sidecar_config_tags`中

**Q: 为什么不用一个item存储value和tags以减小查询的次数，例如:**

```json
{
    "value":"raw value",
    "tags":{
        "tag1":"tag1value",
        "tag2":"tag2value"
    }
}
```

A: 如果采用这种设计方案，那么原来使用apollo的历史遗留项目将不能够迁移到sidecar中

### 如何迁移历史遗留项目

1. Get/subscribe API是兼容的。如果没有使用save/delete的API，那么用户可以很容易地将历史遗留项目将迁移到sidecar当中。只需要在config.json文件中保持`label`字段为空，sidecar将会使用原始key来替代`key@$label`与apollo进行交互。
2. Save/delete API可能不兼容。sidecar会在config.json使用固定的`cluster`和`env`字段。这意味着用户如果想要用其他的appid来变更一些配置项，将无法通过`cluster`和`env`字段作为save/deleteAPI的参数。

### sidecar的config.json文件

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

### 我该使用哪个Apollo SDK？
目前采用官方维护的[SDK](https://github.com/apolloconfig/agollo)，其它Apollo第三方的SDK可以从[链接](https://www.apolloconfig.com/#/zh/usage/third-party-sdks-user-guide)找到。

使用该sdk的一些注意事项:
1. 在连接server和创建client之前，用户必须先在AppConfig声明所有的命名空间，例子如下所示：

```go
	c := &config.AppConfig{
		AppID:          "testApplication_yang",
		Cluster:        "dev",
		IP:             "http://106.54.227.205:8080",
		NamespaceName:  "dubbo",  // 在连接server之前声明
		IsBackupConfig: true,
		Secret:         "6ce3ff7e96a24335a9634fe9abca6d51",
	}
	client, err := agolloConfig.StartWithConfig(func() (*config.AppConfig, error) {
		return c, nil
	})
```

1. 配置不了`env`字段。
2. 无Save/delete API。
3. 无批量查询API。

4. 不能确保并发安全。

5. 配置或者使用不了[Apollo Meta Server](https://www.apolloconfig.com/#/zh/usage/java-sdk-user-guide?id=_122-apollo-meta-server)。

6. 设置kv和tags操作不是事务，存在一致性问题。



### 我应该如何使用Apollo SDK API

| Configuration API      | apollo sdk API                                                                                                                                                                                                                                                                                                 |
| ---------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| GetConfiguration       | cache := client.GetConfigCache(c.NamespaceName)  <br> value,_ := client.Get("key")                                                                                                                                                                                                                             |
| SaveConfiguration      | 通过http. [update](https://www.apolloconfig.com/#/zh/usage/apollo-open-api-platform?id=_3211-%e4%bf%ae%e6%94%b9%e9%85%8d%e7%bd%ae%e6%8e%a5%e5%8f%a3) + [commit](https://www.apolloconfig.com/#/zh/usage/apollo-open-api-platform?id=_3213-%e5%8f%91%e5%b8%83%e9%85%8d%e7%bd%ae%e6%8e%a5%e5%8f%a3) 使用open API |
| DeleteConfiguration    | 使用http. [delete](https://www.apolloconfig.com/#/zh/usage/apollo-open-api-platform?id=_3212-%e5%88%a0%e9%99%a4%e9%85%8d%e7%bd%ae%e6%8e%a5%e5%8f%a3) + [commit](https://www.apolloconfig.com/#/zh/usage/apollo-open-api-platform?id=_3213-%e5%8f%91%e5%b8%83%e9%85%8d%e7%bd%ae%e6%8e%a5%e5%8f%a3) 使用open API |
| SubscribeConfiguration | https://github.com/apolloconfig/agollo/wiki/%E7%9B%91%E5%90%AC%E5%8F%98%E6%9B%B4%E4%BA%8B%E4%BB%B6                                                                                                                                                                                                             |

### 如何订阅指定应用的所有配置变更

订阅config.json文件中所有声明的namespaces。