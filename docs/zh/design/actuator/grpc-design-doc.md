# MOSN gRPC 框架设计文档

## 背景

MOSN 基于 go grpc server 框架提供一个 GRPC Server 的能力，相比于原生的 go grpc server 框架，可以获得如下能力：

+ 完全复用 MOSN Sidecar 的部署、升级、运维能力
+ 可复用 MOSN 中通用的基础能力：热升级、Listener Filter、部分 Network Filter 等
+ 可复用部分 MOSN 中的 StreamFilter 扩展能力

## 设计思路

+ MOSN 的 gRPC 能力主要还是基于官方的 gRPC 库进行核心能力的实现，并且尽量减少 gRPC Server 开发者感受到的差异
+ 基于 NetworkFilter 机制进行实现

## 详细设计

首先梳理一下 NetworkFilter 机制与处理流程

![networkfilter.png](../../../img/actuator/networkfilter.jpg)

+ 在配置解析时，完成 gRPC Server 的启动，随着 MOSN 的 Listener 监听开始提供服务
+ 一个连接对应一个 NetworkFilter 对象
  + InitializeReadFilterCallbacks 和 OnNewConnection 也是在连接创建时调用的接口，负责进行连接初始化的工作
  + OnData 是在收到数据以后调用的接口，负责数据的传递
+ go gRPC Server 库，从“监听”的 Listener 中 Accept 一个连接，然后进行读写，而在 MOSN 框架中，Listener 的监听和连接数据的读写都处理过了，这里需要进行一层处理
  + Listener 和 Conn 都是 interface，可以在 MOSN 的 Filter 中进行处理以后，再把数据传给 gRPC Server，做到 gRPC Server 无感知
  + 在配置解析时，实现 Listener 的封装
  + 在 InitializeReadFilterCallbacks 中实现 Conn 的封装
  + 在 OnNewConnection 中将封装的 Conn 传递给封装的 Listener，触发 Listener.Accept
  + 在 OnData 中将读取到的数据传递给封装的 Conn，触发 Conn.Read

![networkfilter-grpc.png](../../../img/actuator/networkfilter-grpc.jpg)

+ gRPC Server 的实现
  + 在使用官方 gRPC 框架实现 gRPC Server 的时候，开发者需要基于 proto 文件生成一个.pb.go 文件，同时需要实现一组接口满足 proto 中定义的接口实现，将其注册（Register）到 gRPC Server 框架中
  + MOSN 的 gRPC NetworkFilter 也需要提供类似的注册能力，让开发者只关注对应 gRPC Server 实现逻辑，然后注册到 MOSN 框架中即可

+ MOSN GRPC 框架要求开发者实现一个函数，该函数返回一个未调用 Serve 方法的 grpc server。框架会使用自定义的 Listener 去调用 Serve 方法实现对数据的拦截处理

```Go
func init() {
    mgrpc.RegisterServerHandler("mygrpc", MyFunc)
}
func MyFunc(_ json.RawMessage) *grpc.Server {
    s := grpc.NewServer()
    // pb 是.pb.go 所在的 package 路径
    // server 是开发者实现的 api 接口
    pb.RegisterGreeterServer(s, &server{})
    return s
}
```

+ 预期使用的配置示例 (layotto)

```json
{
	"servers":[
		{
			"default_log_path":"stdout",
			"default_log_level": "INFO",
			"listeners":[
				{
					"name":"grpc",
					"address": "0.0.0.0:34904",
					"bind_port": true,
					"filter_chains": [{
						"filters": [
							{
								"type": "grpc",
								"config": {
									"server_name":"runtime",
									"grpc_config": {
										"hellos": {
											"helloworld": {
												"hello": "greeting"
											}
										},
										"config_stores": {
											"etcd": {
												"address": ["127.0.0.1:2379"],
												"timeout": "10"
											}
										}
									}
								}
							}
						]
					}],
					"stream_filters": [
						{
							"type": "flowControlFilter",
							"config": {
								"global_switch": true,
								"limit_key_type": "PATH",
								"rules": [
									{
										"resource": "/spec.proto.runtime.v1.Runtime/SayHello",
										"grade": 1,
										"threshold": 5
									}
								]
							}
						}
					]
				}
			]
		}
	]
}

```
