# 分布式锁组件

**通用配置项说明**

| 字段 | 必填 | 说明 |
| --- | --- | --- |
| keyPrefix | N | key 的前缀策略 |


keyPrefix支持以下键前缀策略:

* **`appid`** - 这是默认策略。用户传入的resource_id最终将被保存为`当前appid||resource_id`

* **`name`** - 此设置使用组件名称作为前缀。 比如redis组件会将用户传入的resource_id存储为`redis||resource_id`

* **`none`** - 此设置不使用前缀。 

*  其他任意不含||的字符串.比如keyPrefix配置成"abc",那么用户传入的resource_id最终将被保存为`abc||resource_id`


**其他配置项**

除了以上通用配置项，每个分布式锁组件有自己的特殊配置项，请参考每个组件的说明文档。