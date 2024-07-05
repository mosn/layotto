# MOSN gRPC framework design document

## Background

MOSN is able to provide a GRPC Server based on the go grpc server framework. Compared to the original go grpc server framework, the following capacity： is available

- Full reuse of MOSN Sidecar deployments, upgrades, shipping capabilities
- Reuse basic：HotUps, Listener Filter, Network Filter, etc. common in MOSN
- StreamFilter extension capability in MOSN

## Idea

- MOSN's gRPC capacity is largely based on the achievement of core competencies in the official gRPC library and varies from gRPC Server developers
- Based on NetworkFilter

## Detailed Design

Split NetworkFilter Mechanisms & Processes

![networkfilter.png](/img/actuator/networkfilter.jpg)

- When configuring the resolution, complete the start of the gRPC Server and start providing the service with the Listener listener listener for MOSN
- A connection to a NetworkFilter object
  - InitializeFilterCallbacks and OnNewConnections are also the interfaces called when the connection is created. They are responsible for initializing connections
  - OnData is an interface called after data is received and is responsible for data transmission
- go gRPC Server library, one connection from Listener to 'listen' and then read and write, while Listener's listening and connection data have been processed in the MOSN framework, where a layer of treatment is required
  - Listener and Conn are interface. Once processed in the MOSN Filter, then get data to gRPC Server, do gRPC Server without knowledge
  - Implementing Listener's package when configuring parsing
  - Implementing Conn packages in InitializeReadFilterCallbacks
  - Transmit the packed Conn to the encapsulated Listener, triggering Listener.Accept
  - Send read data in OnData to encapsulated Conns, Trigger Conn.Read

![networkfilter-grpc.png](/img/actuator/networkfilter-grpc.jpg)

- gRPC Server implementation
  - In implementing the gRPC Server using the official gRPC framework, the developer needs to generate a .pb.go file based on proto file, and a set of interfaces that meet the performance defined in the interface to register (Register) into the gRPC Server framework
  - MOSN gRPC NetworkFilter needs to provide a similar registration capability, allowing developers to focus only on gRPC Server for logic, and then register in MOSN framework enough

- The MOSN GRPC framework requires the developer to implement a function, which returns a grpc server that does not call the Serve method.The framework will use the custom Listener to call the Serve method to block data

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

- Example of intended configuration (layotto)

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
                                                "type": "helloworld",
												"hello": "greeting"
											}
										},
										"config_store": {
											"config_store_demo": {
                                                "type": "etcd",
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
