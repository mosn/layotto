# Nacos

## 配置项说明

示例：configs/config_nacos.json

| 字段                    | 必填  | 说明                                                                                                                              |
|-----------------------|-----|---------------------------------------------------------------------------------------------------------------------------------|
| address               | Y   | nacos 服务地址，填写 <ip:port> 模式即可。默认使用 `/nacos` 作为 url context，使用 http 作为连接协议。如果配置文件中没有填写 acm 连接字段，则这个字段是必填字段。可以填写多个 address 地址进行连接。 |
| timeout               | N   | 连接 nacos 服务的超时时间，单位秒。默认 10s.                                                                                                    |
| metadata.namespace_id | N   | nacos 命名空间，用于配置文件在命名空间级别的隔离。默认为空，代表使用 nacos 的 default 命名空间。                                                                     |
| metadata.username     | N   | nacos 服务 auth 校验所需的用户名。                                                                                                         |
| metadata.password     | N   | nacos 服务 auth 校验所需的密码。                                                                                                          |
| metadata.end_point    | N   | acm 模式字段，表示连接的nacos服务地址，[参考文档](https://help.aliyun.com/document_detail/130146.html)                                             |
| metadata.region_id    | N   | acm 模式字段，表示 nacos 服务所在区域。                                                                                                       |
| metadata.access_key   | N   | acm 模式字段，表示阿里云 nacos 服务中的 access key。                                                                                           |
| metadata.secret_key   | N   | acm 模式字段，表示阿里云 nacos 服务中的 secret key。                                                                                           |
| metadata.log_dir      | N   | nacos go sdk 输出日志文件目录地址。默认目录名称为`/tmp/layotto/nacos/logs`，默认日志名称为`nacos-sdk.log`。                                                    |
| metadata.log_level    | N   | nacos go sdk 输出日志等级。支持日志等级为 `debug`,`info`,`warn`,`error`。默认日志等级为`debug`。                                                                |
| metadata.cache_dir    | N   | nacos 配置本地缓存文件路径。默认缓存路径为`/tmp/layotto/nacos/cache`。                                                                               |

> 当有任一 acm 字段出现在配置文件中，则 layotto 采用 acm 模式连接到 nacos 服务，忽略 address、username、password 等字段。

## 其他配置项

还需要配置 `app_id` 字段，是必填字段。表示该 app 的类别。支持在配置中心层面对不同的 app 服务进行配置隔离。

不过不是添加在 nacos configstore 组件的配置里，添加在而外的配置项中，方便其他需要使用 `app_id` 的组件复用。

![img.png](../../../img/configuration/nacos/img.png)

## 怎么启动 Nacos

nacos 的启动方式可以参考[nacos 官方文档](https://nacos.io/zh-cn/docs/quick-start-docker.html)

部署后需要修改Layotto的[config文件](https://github.com/mosn/layotto/blob/main/configs/config_nacos.json) ，将nacos服务器地址等信息改成您自己的。
