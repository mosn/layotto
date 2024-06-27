# 状态管理组件
**配置文件结构**

json配置文件有如下结构：

```json
"state": {
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
"state": {
  "state_demo": {
    "type": "redis",
    "metadata": {
      "redisHost": "localhost:6380",
      "redisPassword": ""
    }
  }
}
```


**通用配置项说明**

不同State组件的通用配置项有：

| 字段 | 必填 | 说明 |
| --- | --- | --- |
| keyPrefix | N | key 的前缀策略 |


keyPrefix支持以下键前缀策略:

* **`appid`** - 用户传入的key最终将被保存为`当前appid||key`

* **`name`** - 此设置使用组件名称作为前缀。 比如redis组件会将用户传入的key存储为`redis||key`

* **`none`** - 不给 key 添加前缀。**这是默认策略。**

*  其他任意不含||的字符串.比如keyPrefix配置成"abc",那么用户传入的key最终将被保存为`abc||key`


**其他配置项**

除了以上通用配置项，每个State组件有自己的特殊配置项，请参考每个组件的说明文档。