# 分布式（自增）唯一id组件
**配置文件结构**

json配置文件有如下结构：

```json
"sequencer": {
  "<Component A Name>": {
    "type": "Component A Name",
    "biggerThan": {
      "<KEY>": "<VALUE>",
      "<KEY>": "<VALUE>"
    },
    "metadata": {
      "<KEY>": "<VALUE>",
      "<KEY>": "<VALUE>"
    }
  },
  "<Component B Name>": {
    "type": "Component B Name",
    "biggerThan": {
      "<KEY>": "<VALUE>",
      "<KEY>": "<VALUE>"
    },
    "metadata": {
      "<KEY>": "<VALUE>",
      "<KEY>": "<VALUE>"
    }
  }
}
```

您可以在metadata里配置组件关心的key/value配置。例如[etcd组件的配置](https://github.com/mosn/layotto/blob/main/configs/runtime_config.json) 如下：

```json
"sequencer": {
  "sequencer_demo": {
    "type": "etcd",
    "biggerThan": {
      "key1": 1,
      "key2": 111
    },
    "metadata": {
      "endpoints": "localhost:2379",
      "segmentCacheEnable": "false",
      "segmentStep": "1",
      "username": "",
      "password": "",
      "dialTimeout": "5"
    }
  }
},
```

**通用配置项说明**

| 字段 | 必填 | 说明 |
| --- | --- | --- |
| biggerThan | N | 要求组件生成的所有id都得比"biggerThan"大。设计这个配置项是为了方便用户做移植。比如系统原先使用mysql做发号服务，id已经生成到了1000，后来迁移到PostgreSQL上，需要配置biggerThan为1000，这样PostgreSQL组件在初始化的时候会进行设置、强制id在1000以上,或者发现id没法满足要求、直接启动时报错。 |
| segmentCacheEnable | N | 是否开启号段缓存。默认值true |
| segmentStep | N | 每次号段缓存的大小，默认值50 |

- 什么是segment(号段)模式?

原始方案是每次获取ID都得读写一次数据库，造成数据库压力大。segment(号段)模式的设计目的是想在sidecar里面提前缓存一些id（缓存一个号段内的所有id),为数据库减压。

具体来说，Layotto每次从sequencer组件批量获取id，即每次获取一个segment(segmentStep决定大小)号段的值。用完之后Layotto再调用组件、组件再去数据库获取新的号段，可以大大的减轻数据库的压力。

这种设计参考了[美团Leaf的设计](https://tech.meituan.com/2017/04/21/mt-leaf.html)

**其他配置项**

除了以上通用配置项，每个组件有自己的特殊配置项，请参考每个组件的说明文档。