# 发布/订阅组件
**配置文件结构**

json配置文件有如下结构：

```json
"pub_subs": {
  "<Component A Name>": {
    "type":"<Component A Type>",
    "metadata": {
      "<KEY>": "<VALUE>",
      "<KEY>": "<VALUE>"
    }
  },
  "<Component B Name>": {
    "type":"<Component B Type>",
    "metadata": {
      "<KEY>": "<VALUE>",
      "<KEY>": "<VALUE>"
    }
  }
}
```

您可以在metadata里配置组件关心的key/value配置。例如[redis组件的配置](https://github.com/mosn/layotto/blob/main/configs/config_redis.json) 如下：

```json
"pub_subs": {
  "pubsub_dmoe": {
    "type": "redis",
    "metadata": {
      "redisHost": "localhost:6380",
      "redisPassword": ""
    }
  }
},
```


**配置项说明**

每个State组件有自己的特殊配置项，请参考每个组件的说明文档。