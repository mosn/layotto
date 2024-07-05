# Layotto go-sdk 使用指南 

## 从 HelloWorld 入手

### 快读体验一下 go-sdk 使用

```shell
# 1. 在项目根目录下按序执行，创建 layotto 二进制程序，加入到环境变量中方便后续操作
cd cmd/layotto
go build -o layotto
export PATH=$PATH:$(pwd)/layotto

# 2.根据配置文件启动 layotto.
layotto start -c ../../configs/config_standalone.json

# 3.另起一个终端，启动 go-sdk 客户端程序
cd demo/hello/common
go build -o client
./client -s helloworld

## 输出结果 output:
runtime client initializing for: 127.0.0.1:34904
greeting
```

### go-sdk 使用代码

简单观察一下 `client` 的使用流程，代码放在 `demo/hello/common/client.go` 中，下面是缩略版的内容，内容如下

```go
var storeName string

func init() {
	flag.StringVar(&storeName, "s", "", "set `storeName`")
}

func main() {
    flag.Parse()
    if storeName == "" {
		panic("storeName is empty.")
	}
    
    // 创建连接，采用默认IP端口 127.0.0.1:34904，也是本地 layotto 的默认启动端口
	cli, err := client.NewClient()
	if err != nil {
		log.Fatal(err)
	}

    // 关闭 grpc 连接
	defer cli.Close()	
	
    // 调用 go-sdk 接口，发送请求
	res, err := cli.SayHello(context.Background(), &client.SayHelloRequest{
		ServiceName: storeName,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(res.Hello)
}
```

### 相关配置文件

可以通过使用一个最简单配置文件来启动 Layotto,测试上面的 HelloWorld 程序. 内容如下：

> 有关配置文件更详细的介绍可以参考 [Layotto 配置文件介绍](zh/configuration/overview.md).

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


在上述流程中，Layotto 注册了配置文件中 `grpc_config` 属性里的组件.

组件的实现的接口类型是 `hellos`，下面可以支持多个 instance 实例实现这个接口，配置文件中的 instance 属性是用于区分不同实现该接口的实例名称. 上述配置文件中，只涉及到一个 instance `helloworld`.

在使用 go-sdk 调用 hello 相关的服务，例如程序中的 `SayHello`，需要指定 instance 的名称，上述代码使用 `ServiceName` 参数传递 instance 名称，表明使用指定 instance 来操作这次请求.

instance 属性里还有 `type` 和 `hello` 两个子属性：

- `type` 属性是每个 instance 都有的，表明使用什么类型来实现了该接口，这里是 `helloworld` 类型. 这个 `type` 目前要求是已经注册到 Layotto 中组件类型，可以通过 [cmd/layotto/main.go](https://github.com/mosn/layotto/blob/6f6508b11783f1e4fa947ff47632e74064333384/cmd/layotto/main.go#LL275C1-L275C1) 中查看注册的组件类型.
也可以查看各种类型组件配置文件的[文档](zh/configuration/overview.md). 

- `hello` 属性用于初始化组件，不同类型的接口有不同的属性字段. 不同 `type` 的 instance 也有不同的属性字段来初始化. 可以通过[代码](https://github.com/mosn/layotto/blob/6f6508b11783f1e4fa947ff47632e74064333384/components/hello/hello.go#L32-L36)查看相关配置属性的用法.

## 初始化实例

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


## 体验其他接口函数

以分布式id `sequencer` 为例，使用 `redis` 类型，相关参数请参考这篇[文档](zh/component_specs/sequencer/redis.md).

具体代码配置参考 `demo/sequencer/common/client.go`，内容如下.

```go
const (
    key = "key666"
)

var storeName string

func init() {
    flag.StringVar(&storeName, "s", "", "set `storeName`")
}

func main() {
    cli, err := client.NewClient()
	if err != nil {
		panic(err)
	}
	defer cli.Close()
	ctx := context.Background()
	fmt.Printf("Try to get next id.Key:%s \n", key)
	for i := 0; i < 10; i++ {
        
        // 获取下一个分布式id
		id, err := cli.GetNextId(ctx, &runtimev1pb.GetNextIdRequest{
			StoreName: storeName, // 这里是 storename 就是上文中提到的 instance 名称, 需要与配置文件中填写的 instance 保持一致.
			Key:       key,
			Options:   nil,
			Metadata:  nil,
		})
		if err != nil {
			panic(err)
		}
		fmt.Printf("Next id:%v \n", id)
	}
	fmt.Println("Demo success!")
}
```

在使用 `GetNextId` 函数接口时传递了一下几个参数：
- StoreName： 表明处理此次请求的 instance 名称，与上文中 `SayHello` 接口中 `ServiceName` 的作用一致. 
- Key： 分布式 id 的 namespace，不同的 namespace 允许使用相同的分布式 id.
- Options： 填写id生成的模式，有 strong 和 weak 可供选择. 前者生成的id是绝对比前一次请求的 id 大，但性能上会有所损耗. 后者在性能上有优势，但是无法满足 id 绝对比前一个大. 默认使用 weak 模式生成 id.
- Metadata： 其他参数. 不同类型的 instance 需要不同的 metadata 进行定制化构建. 例如，[mysql 类型](https://github.com/mosn/layotto/blob/6f6508b11783f1e4fa947ff47632e74064333384/components/pkg/utils/mysql.go#L42-L60)的 sequencer 支持使用 metadata 传递 `tableName`，`defaultPassword` `userName` 等信息，而 redis 类型的 instance 则不需要传递 metadata.

返回：
- next_id： 分布式id.


> 关于接口的详细使用方式，可以通过查看[API设计文档](https://mosn.io/layotto/#/zh/design/sequencer/design). 可能文档还未即使补充，也可以查看[相关代码注释](https://github.com/mosn/layotto/blob/6f6508b11783f1e4fa947ff47632e74064333384/spec/proto/runtime/v1/runtime.pb.go#L1118-L1131)，及[组件内部实现](https://github.com/mosn/layotto/blob/6f6508b11783f1e4fa947ff47632e74064333384/pkg/grpc/default_api/api_sequencer.go#L99-L119)获悉.



相关参数配置参数如下：

```json
"sequencer": {
    "sequencer_demo": {
        "type": "redis",
        "metadata": {
            "redisHost": "127.0.0.1:6380",
            "redisPassword": ""
        }
    }
}
```

启动应用程序进行验证.

```shell
# docker 启动 redis
docker run --name layotto-redis -itd -p 6380:6379 redis 

# 启动 layotto
cd configs
layotto start -c config_redis.json 

# 另起一个终端，开启程序
cd demo/sequencer/common
go build -o client
./client -s sequencer_demo

# output
runtime client initializing for: 127.0.0.1:34904
Try to get next id.Key:key666
Next id:next_id:1
Next id:next_id:2
Next id:next_id:3
Next id:next_id:4
Next id:next_id:5
Next id:next_id:6
Next id:next_id:7
Next id:next_id:8
Next id:next_id:9
Next id:next_id:10
Demo success!


# 进入 redis 容器，查看相关修改。
docker exec -it layotto-redis redis-cli -h localhost -p 6379
localhost:6379> keys *
1) "sequencer|||app1||key666"
localhost:6379> get "sequencer|||app1||key666"
"10000"

```



## 原生 GRPC 连接调用接口 `go-sdk` 未包含接口 

目前，go-sdk 只是对 grpc 做了一层很薄的封装，所以对于使用未在 go-sdk 同步的接口，可以采用 grpc 的方式进行调用.

### runtimeAPI 接口

这里以 `file` 接口中的 `local` 组件类型为例。

通过可以查看 [`spec/proto/runtime`](https://github.com/mosn/layotto/tree/main/spec/proto/runtime/v1) 中的 proto 文件或[GRPC API文档](zh/api_reference/README.md)

完整代码参考 `demo/file/local/client.go`，内容如下：

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


在 local 模式下的 `Put` 创建需要设置 Metadata 中的 `FileMode` 和 `FileFlag` 参数，来创建文件。

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

## extension 接口

目前 go-sdk 有部分**建设中**接口如 sms，cryption，email 等并未放在 `runtimeAPI` service 中，而是使用 proto 生成的 interface 接口. 所以在使用此类接口时，可以像 file 接口直接调用 grpc 接口接口. 相关接口文档参考 [spec 目录下的 proto 文件](https://github.com/mosn/layotto/tree/main/spec/proto/extension/v1)

下面以 email 接口为例，内容如下：

> 接口还处于建设中，没有 component 实现，所以无法真正调用，下面仅为简单示例。

```go
func main() {
	cli, err := client.NewClient()
	if err != nil {
		panic(err)
	}
	
	defer cli.Close()

	template, err := cli.SendEmailWithTemplate(context.TODO(), &pb.SendEmailWithTemplateRequest{})
	if err != nil {
		fmt.Println(err)
		return 
	}
	
	fmt.Println(template.String())
}
```



## 更多示例

sdk 其他接口使用可以参考 [demo目录下的代码示例](https://github.com/mosn/layotto/tree/main/demo) 及 [quick-start 启动文档](zh/start/README.md)

查看相关配置文件编写可以参考 [configs示例](https://github.com/mosn/layotto/tree/main/configs) 及 [组件配置文档](zh/configuration/overview.md) 

