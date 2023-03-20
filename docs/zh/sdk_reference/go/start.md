# Layotto go-sdk 使用指南 

## 从 HelloWorld 入手

简单观察一下 `client` 的使用流程，代码在 `sdk/go-sdk/demo/hello/client.go` 中

```go
func main() {
    // 创建连接，采用默认IP端口 127.0.0.1:34904
	cli, err := client.NewClient()
	if err != nil {
		log.Fatal(err)
	}

    // 关闭 grpc 连接
	defer cli.Close()	
	
    // 调用 go-sdk 接口，发送请求
	res, err := cli.SayHello(context.Background(), &client.SayHelloRequest{
		ServiceName: "helloworld", // 这个是 hellos 的 instance，与配置文件中需设置一致 
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(res.Hello)
}
```

**关于 `Client` 支持的函数接口可以查看 [`sdk/go-sdk/client/client.go`](https://github.com/mosn/layotto/blob/main/sdk/go-sdk/client/client.go) 中的 interface** 

**目前 `go-sdk` 支持的初始化方式**

```go
// 使用默认的端口 "34904"，也可以用 runtime_GRPC_PORT 环境变量设置端口，同时使用默认的 IP 地址 "127.0.0.1"，内部调用 NewClientWithPort 函数
NewClient() (client Client, err error)

// 可以设置端口，使用默认的 IP 地址 "127.0.0.1"，内部调用 NewClientWithAddress 函数
NewClientWithPort(port string) (client Client, err error)

// 要求传入 Address, 即 <ip:port> 的形式，创建完 GRPC conn 后调用 NewClientWithConnection
NewClientWithAddress(address string) (client Client, err error)

// 传入用户自己创建好的 GRPC 连接封装成 Client
NewClientWithConnection(conn *grpc.ClientConn) Client
```



使用一个**最简单配置文件**来启动 Layotto 测试我们的 `HelloWorld` 程序。我放在了 `sdk/go-sdk/demo/hello/config.json` 中。

```json
{
  "servers": [
    {
      "default_log_path": "stdout",
      "default_log_level": "DEBUG",
      "listeners": [
        {
          "name": "grpc",
          "address": "127.0.0.1:34904",
          "bind_port": true,
          "filter_chains": [
            {
              "filters": [
                {
                  "type": "grpc",
                  "config": {
                    "server_name": "runtime",
                    "grpc_config": {
                      "hellos": {
                        "helloworld": {
                          "type": "helloworld",
                          "hello": "greeting"
                        }
                      }
                    }
                  }
                }
              ]
            }
          ]
        }
      ]
    }
  ]
}
```



**简单体验一下这个 Demo。**

```shell
## 在项目根目录下按序执行
## 创建 layotto 二进制程序
cd cmd/layotto
go build -o layotto
## 根据配置文件启动 layotto, linux 下可以用 export 添加到环境变量中或拷贝至 /usr/local/bin 目录下
./layotto start -c ../../sdk/go-sdk/demo/hello/config.json

## 另起一个终端(项目根目录下)，启动我们的客户端
cd sdk/go-sdk/demo/hello/
go build -o client
./client

## 输出结果 output:
runtime client initializing for: 127.0.0.1:34904
greeting
```



## 关于配置文件

Layotto 本质是在 Mosn 中开启一个 GRPC 服务，监听指定的端口供客户端程序连接。所以在配置文件中，我们需要配置一个 Listener，其他的配置需求可以查看 Mosn 的[相关文档](https://mosn.io/docs/products/configuration-overview/)。

### 组件名

同时，Layotto 支持多个**同类型**的组件，如针对 `hellos` 这个接口，用户可以设置多个 instance 实例，`hellos` 这个名称是强要求的，可以参考`pkg/runtime/config.go` 中 `MosnRuntimeConfig` 结构体中对应的类型；

```go
type MosnRuntimeConfig struct {
	AppManagement          AppConfig                           `json:"app"`
	HelloServiceManagement map[string]hello.HelloConfig        `json:"hellos"`
	ConfigStoreManagement  map[string]configstores.StoreConfig `json:"config_store"`
	RpcManagement          map[string]rpc.RpcConfig            `json:"rpcs"`
	PubSubManagement       map[string]pubsub.Config            `json:"pub_subs"`
	StateManagement        map[string]state.Config             `json:"state"`
	Files                  map[string]file.FileConfig          `json:"file"`
	Oss                    map[string]oss.Config               `json:"oss"`
	LockManagement         map[string]lock.Config              `json:"lock"`
	SequencerManagement    map[string]sequencer.Config         `json:"sequencer"`
	Bindings               map[string]bindings.Metadata        `json:"bindings"`
	SecretStoresManagement map[string]secretstores.Metadata    `json:"secret_store"`
	// <component type,component name,config>
	// e.g. <"super_pubsub","etcd",config>
	CustomComponent map[string]map[string]custom.Config `json:"custom_component,omitempty"`
	Extends         map[string]json.RawMessage          `json:"extends,omitempty"` // extend config
	ExtensionComponentConfig
}
```

### Instance 命名

而实例的名称是不做要求的，可以随意命名。主要用于客户端在调用相关接口时传入指定的 instance 名称，要求该 instance 来完成接口做的工作。

拿上面的配置文件与程序 `demo` 为例，我们在 `hellos` 下设置了 `helloworld` 实例，并在程序中的 `ServiceName`  中填入该 `instance`  的名称。如果我们 `ServiceName` 与配置文件中填写的**不一致**，如从 "hellowrold" 改成 "hello_demo"，会遇到下面这个**错误**。

```shell
runtime client initializing for: 127.0.0.1:34904
2023/03/19 22:33:28 rpc error: code = Unknown desc = no instance found
```

### Instance 配置

**Type**  类型可以参考 [`cmd/layotto/main.go`](https://github.com/mosn/layotto/blob/main/cmd/layotto/main.go) 中 `runtime` 注册的名称。`type` 字段未填或填错，会遇到一下错误。通过配置错误中的 `instance-name` 排查一下配置文件。

```shell
# [Info] ...
2023-03-20 10:26:22,646 [ERROR] [runtime] occurs an error: service component  is not regsitered, create sequencer component <instance-name> failed       
2023-03-20 10:26:22,646 [ERROR] create a registered server failed: service component  is not regsitered
# ...
panic: service component  is not regsitered
# ....
```



 **config** 结构体我们需要参考 `component` 中对应的 `config` 结构体
以 `hellos` 组件为例，我们可以从 `component/hello/hello.go` 找到这个结构体。

```go
type HelloConfig struct {
	ref.Config
	Type        string `json:"type"`
	HelloString string `json:"hello"`
}
```

**对于其他 instance 配置编写也是如此。**



##  尝试其他 API 接口

> 知道了怎么配置配置文件和调用接口，现在开始自己试试吧。

我们以 `sequencer` 中的 `in-memory` 类型为例。先寻找他的配置结构体。

```go
type Config struct {
	ref.Config
	Type       string            `json:"type"`
	BiggerThan map[string]int64  `json:"biggerThan"`
	Metadata   map[string]string `json:"metadata"`
}
```

但是通过查看 `in-memory` 的 `Init` 函数，发现他其实不需要使用 `Config` 中的参数。

```go
func (s *InMemorySequencer) Init(_ sequencer.Configuration) error {
	readinessIndicator.SetStarted()
	livenessIndicator.SetStarted()
	return nil
}
```

所以在配置文件中我们只需要填写 `Type` 字段。

```json
"sequencer": {
    "sequencer-demo": {
      "type": "in-memory"
    }
}
```

**使用 `go-sdk` 调用 sequencer 接口** 

```go
var (
	storeName = "sequencer-demo"
)

func main() {
	cli, err := client.NewClient()
	if err != nil {
		fmt.Println("connect to layotto failed")
		panic(err)
	}
    
    defer cli.Close()

    // 通过观察，发现 in-memory sequencer 中只需要使用到 key 字段。
	id, err := cli.GetNextId(context.Background(), &runtimev1pb.GetNextIdRequest{
		StoreName: storeName,
		Key:       "key",
	})

	if err != nil {
		fmt.Println("get next id failed")
		panic(err)
	}
	fmt.Println(id)
}

/*
output:
runtime client initializing for: 127.0.0.1:34904
next_id:10
*/
```

**具体代码配置参考`sdk/go-sdk/demo/sequencer/in-memory`**



## 针对同一批接口切换 instance 实例

我们使用 redis 存储 sequencer 来替代 in-memory 类型。

```shell
cd sdk/go-sdk/demo/sequencer/redis
## docker-compose 启动 redis 和 layotto
docker run --name layotto-redis -itd -p 36379:6379 redis 
layotto start -c config.json
```

配置文件相关参数，只填写了必要的字段。
```json
"sequencer": {
    "sequencer-demo": {
        "type": "redis",
        "metadata": {
        "redisHost": "127.0.0.1:36379"
        }
    }
}
```

程序代码与 `in-memeory` 的一致。运行程序，查看日志输出。

```shell
go build -o client
./client

# output
runtime client initializing for: 127.0.0.1:34904
next_id:10
```

进入 redis 容器，查看相关修改。

```shell
docker exec -it layotto-redis redis-cli -h localhost -p 6379
localhost:6379> keys *
1) "sequencer|||key"
localhost:6379> get  "sequencer|||key"
"10000"
```

**具体代码配置参考`sdk/go-sdk/demo/sequencer/redis`**


## 原生 GRPC 连接调用接口 `go-sdk` 未包含接口 

### runtimeAPI 接口

我们可以查看 [`spec/proto/runtime`](https://github.com/mosn/layotto/tree/main/spec/proto/runtime/v1) 中的 proto 文件，也可以查看[GRPC API文档](zh/api_reference/README.md)

比如 `go-sdk` 目前还未支持 `file`  接口，我们可以使用 `proto`文件中提供的即可来调用。 这里以 `file` 中的 `local` Type 为例。

```go
// 创建 grpc 连接，包装 runtimeclient
func main() {
	// conn to layotto grpc server with row grpc client
	conn, err := grpc.Dial("127.0.0.1:34904", grpc.WithInsecure())
	if err != nil {
		fmt.Printf("conn build failed,err:%+v", err)
		return
	}

 	c := runtimev1pb.NewRuntimeClient(conn)
    
    // 调用 file 的相关接口
    // ....
}
```



在调用接口时传递的**参数**需要参考对应 `component` 中对应的函数。如 local 模式下的 `Put` 创建需要设置 Metadata 中的 `FileMode` 和 `FileFlag` 参数，来创建文件。不同的 Instance 需要传递的参数不一样。

```go
// layotto server
func (lf *LocalStore) Put(ctx context.Context, f *file.PutFileStu) error {
	// 这里需要 f 中的 Metadata 设置 FileMode 和 FileFlag 参数
    mode, ok := f.Metadata[FileMode]
	if !ok {
		return fmt.Errorf("fileMode is required for put file")
	}
	flag, ok := f.Metadata[FileFlag]
	if !ok {
		return fmt.Errorf("fileFlag is required for put file")
	}
	// ....
}


// client
func main() {
    // grpc conn and set runtime client
    
    // 这里我们就要设置相关的 FileMode 和 FileFlag
	meta := make(map[string]string)
	meta["FileMode"] = "0644"
	meta["FileFlag"] = strconv.Itoa(os.O_CREATE | os.O_RDWR)
	req := &runtimev1pb.PutFileRequest{StoreName: StoreName, Name: FileName, Metadata: meta}
	stream, err := c.PutFile(context.Background())
	// ...
}
```

而具体的 GRPC  客户端的实现需要参考 `pkg/grpc/default_api/api.go` 中函数的实现。

**具体的代码参考 `sdk/go-sdk/demo/file/client.go`**



### 其他 extension 接口

针对后加入的接口如 `email` `phone` `sms` `oss` `delay-queue` 等，`go-sdk` 也还未实现他们的 sdk 封装，如果需要使用他们也需要参考对应的 [grpc 接口](https://github.com/mosn/layotto/tree/main/spec/proto/extension/v1)。

go-sdk 封装了他们的 grpc client 接口类型。所以，要想调用它们，我们还是需要通过 GRPC 的方式进行调度。 

```go
// Client is the interface for runtime client implementation.
type Client interface {
	runtimeAPI

	s3.ObjectStorageServiceClient

	// "mosn.io/layotto/spec/proto/extension/v1/cryption"
	cryption.CryptionServiceClient

	// "mosn.io/layotto/spec/proto/extension/v1/delay_queue"
	delay_queue.DelayQueueClient

	// "mosn.io/layotto/spec/proto/extension/v1/email"
	email.EmailServiceClient

	// "mosn.io/layotto/spec/proto/extension/v1/phone"
	phone.PhoneCallServiceClient

	// "mosn.io/layotto/spec/proto/extension/v1/sms"
	sms.SmsServiceClient
}
```



## `go-sdk` 接口使用文档

> 如果文档还未更新，可以自行阅读 `sdk/go-sdk` 中的源码，欢迎提出错误贡献 sdk 使用文档。
>

#### runtimeAPI

- [hello](zh/sdk_reference/go/hello.md)  
- [Invoke](zh/sdk_reference/go/invoke.md)
- [PublishEvent](zh/sdk_reference/go/publish_event.md)
- [State](zh/sdk_reference/go/state.md)
- [Lock](zh/sdk_reference/go/lock.md)
- [Sequencer](zh/sdk_reference/go/sequencer.md)
- [Secret](zh/sdk_reference/go/secret.md)
- [Configuration](zh/sdk_reference/go/config.md)



#### extension

> 目前该部分 go-sdk 仅支持了使用 grpc 接口的方式调用，还未实现封装，可以参考 [GRPC API 文档](zh/api_reference/README.md) 或参考对应的源代码 [`spec/proto/extension`](https://github.com/mosn/layotto/tree/main/spec/proto/extension/v1) 或 [`demo`](https://github.com/mosn/layotto/tree/main/demo) 中的相应使用案例。
> 

- [ObjectStorage](zh/sdk_reference/go/extension/object_storage.md)
- [Cryption](zh/sdk_reference/go/extension/cryption.md)
- [DelayQueue](zh/sdk_reference/go/extension/delay_queue.md)
- [Phone](zh/sdk_reference/go/extension/phone.md)
- [Email](zh/sdk_reference/go/extension/email.md)
- [Sms](zh/sdk_reference/go/extension/sms.md)
