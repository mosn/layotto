# Layotto 配置文件介绍
Layotto 启动时需要读取一个 json 格式的配置文件。示例配置文件：configs/config_apollo.json

目前，Layotto使用一个MOSN 4层filter与MOSN集成、跑在MOSN上，**所以Layotto用到的配置文件其实就是MOSN配置文件。**

![img.png](../../img/configuration/layotto/img.png)

如上图示例，大部分配置是MOSN的配置项，参考[MOSN的配置说明](https://mosn.io/docs/configuration/) ;

其中`"type":"grpc"`对应的filter是MOSN的一个4层filter，用于把Layotto和MOSN集成到一起。

而`grpc_config`里面的配置项是Layotto的"组件配置"，结构为：

```json
"grpc_config": {
  "<API NAME>": {
    "<COMPONENT A NAME>": {
      "type": "<COMPONENT A Type>"
      "<KEY>": "<VALUE>",
      "metadata": {
        "<KEY>": "<VALUE>",
        "<KEY>": "<VALUE>"
      }
    }
  },
  "<API NAME>": {
    "<COMPONENT B NAME>": {
      "type": "<COMPONENT B Type>"
      "<KEY>": "<VALUE>",
      "metadata": {
        "<KEY>": "<VALUE>",
        "<KEY>": "<VALUE>"
      }
    }
  },
}
```

至于每个API NAME填啥、每个组件名是啥、组件能配哪些Key/Value配置项，您可以查阅[组件配置说明](zh/component_specs/overview)