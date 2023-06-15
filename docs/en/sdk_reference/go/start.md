# Layotto go-sdk Reference

## Start from HelloWorld 

### go-sdk simple usage 

```shell
# 1. Execute the following commands sequentially in the root directory of the project
# build a layotto binary program, and add it to the environment variables for easy subsequent operations
cd cmd/layotto
go build -o layotto
export PATH=$PATH:$(pwd)/layotto

# 2.Run layotto according to the configuration file
layotto start -c ../../configs/config_hello.json

# 3.Start another terminal and launch the go sdk client program
cd demo/hello/common
go build -o client
./client -s helloworld

## output:
runtime client initializing for: 127.0.0.1:34904
greeting
```

### The go-sdk code

Briefly observe the usage process of `client`, and the source code is placed in `demo/hello/common/client.go`

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
    
    // Create a connection using the default IP port 127.0.0.1:34904, which is also the default startup port for local Layotto
	cli, err := client.NewClient()
	if err != nil {
		log.Fatal(err)
	}

    // Close grpc connection
	defer cli.Close()	
	
	// call the interface, send request
	res, err := cli.SayHello(context.Background(), &client.SayHelloRequest{
		ServiceName: storeName,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(res.Hello)
}
```

### Configuration files

Start Layotto and test the HelloWorld program above by using the simplest configuration file, the content is as follows:

> For a more detailed introduction to configuration files, please refer to [configuration](en/configuration/overview.md).

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

In the above process, Layotto registers components in the `grpc_config` attribute in the configuration file.

The interface type of the component implementation is `hellos`, and multiple instance instances can be supported to implement this interface.
The instance attribute in the configuration file is used to distinguish the instance names of different implementations of this interface. 
In the above configuration file, only one instance `helloworld` is involved.


When using go sdk to call hello related services, such as `SayHello` in the program, the name of the instance needs to be specified.
The above code uses the `ServiceName` parameter to pass the instance name, indicating that the specified instance is used to operate this request

There are two sub attributes in the instance attribute `type` and `hello`:
- `type`, this attribute is unique to every instance, indicating which type is used to implement the interface. Here, it is the `hellowrold` type which is already been registered as a component type in [Layotto](https://github.com/mosn/layotto/blob/6f6508b11783f1e4fa947ff47632e74064333384/cmd/layotto/main.go#LL275C1-L275C1) 
- `hello`, this attribute is used to initialize components, and different types of interfaces have different attribute fields. Different instance also have different attribute field to initialize. Use [code](https://github.com/mosn/layotto/blob/6f6508b11783f1e4fa947ff47632e74064333384/components/hello/hello.go#L32-L36) to view the usage of relevant configuration.

## Initialize Client

**The initialization methods currently supported by `go sdk`**

```go
// Use the default port "34904", or you can use "runtime_GRPC_Set" the port for the PORT environment variable and use the default IP address "127.0.0.1" to internally call the NewClientWithPort function
NewClient() (client Client, err error)

// Provide port input parameters, use the default IP address "127.0.0.1", and internally call the NewClientWithAddress function
NewClientWithPort(port string) (client Client, err error)

// The address is required to be passed in, that is, in the form of <ip:port>. After creating the GRPC conn, call NewClientWithConnection
NewClientWithAddress(address string) (client Client, err error)

// Encapsulate the GRPC connection created by the user into a client
NewClientWithConnection(conn *grpc.ClientConn) Client
```


## Experience other interface functions

Taking the distributed id `sequencer` as an example, use the 'redis' type. Please refer to this [document](en/component_specs/sequencer/redis. md) for relevant parameters usage.

Please refer to `demo/sequencer/common/client.go` for specific code configuration, as follows:

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
        
		// get the next sequencer id
		id, err := cli.GetNextId(ctx, &runtimev1pb.GetNextIdRequest{
			// Here is the storename, which is the instance name mentioned earlier and needs to be 
			// consistent with the instance filled in the configuration file
			StoreName: storeName, 
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

Some parameters when using `GetNextId` interface:
- StoreName: Indicates the name of the instance processing this request, which is consistent with the role of `ServiceName` in the `SayHello` interface mentioned earlier.
- Key: Distributed ID namespaces allow different namespaces to use the same distributed ID
- Options: Fill in the mode of ID generation, with **strong** and **weak** options to choose from. The ID generated by the **strong** mode is definitely larger than the ID from the previous request, but there may be some performance loss. The **weak** mode has performance advantages, but cannot meet the requirement that the ID is definitely larger than the previous one. By default, use the weak mode to generate id.
- Metadata: Other parameters. Different types of instances require different metadata for customized construction. For example, [MySQL](https://github.com/mosn/layotto/blob/6f6508b11783f1e4fa947ff47632e74064333384/components/pkg/utils/mysql.go#L42-L60) The sequencer supports using metadata to pass information such as `tableName`, `defaultPassword`, and `userName`, while instances of Redis type do not need to pass metadata
returns:
- next_id: The distributed id.


> For detailed useage of interface, you can refer to the [API design docs](https://mosn.io/layotto/#/en/design/sequencer/design)
> Perhaps the document has not been supplemented yet, you can still view the [relecant code comments](https://github.com/mosn/layotto/blob/6f6508b11783f1e4fa947ff47632e74064333384/spec/proto/runtime/v1/runtime.pb.go#L1118-L1131)

The relevant parameter configuration parameters are as follows:

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

Launch the application for verification

```shell
# use docker to start redis
docker run --name layotto-redis -itd -p 6380:6379 redis 

# start layotto
cd configs
layotto start -c config_redis.json 

# Start another terminal and start the client program
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


# exec into redis container and check the key.
docker exec -it layotto-redis redis-cli -h localhost -p 6379
localhost:6379> keys *
1) "sequencer|||app1||key666"
localhost:6379> get "sequencer|||app1||key666"
"10000"

```

## Experience the GRPC interface

Currently, go-sdk only encapsulates GRPC with a thin layer, so for interfaces that are not synchronized in go sdk, GRPC can be used for calling

### runtimeAPI interface

Here, take the `local` component type in the `file` interface as an example.

You can view the proto files in [`spec/proto/runtime`](https://github.com/mosn/layotto/tree/main/spec/proto/runtime/v1) or the [GRPC API docs](en/api_reference/README.md)

The complete code reference is `demo/file/local/client.go`, and the content is as follows:

```go
// create grpc connection, packaged as a runtimeclient
func main() {
	// conn to layotto grpc server with row grpc client
	conn, err := grpc.Dial("127.0.0.1:34904", grpc.WithInsecure())
	if err != nil {
		fmt.Printf("conn build failed,err:%+v", err)
		return
	}

 	c := runtimev1pb.NewRuntimeClient(conn)
    
    // call file's interface
    // ....
}
```

Creating `Put` in local mode requires setting the `FileMode` and `FileFlag` parameters in metadata to create files.

```go
// layotto server
func (lf *LocalStore) Put(ctx context.Context, f *file.PutFileStu) error {
	// Here, you need to set the FileMode and FileFlag parameters using the metadata
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
   
	meta := make(map[string]string)
	meta["FileMode"] = "0644"
	meta["FileFlag"] = strconv.Itoa(os.O_CREATE | os.O_RDWR)
	req := &runtimev1pb.PutFileRequest{StoreName: StoreName, Name: FileName, Metadata: meta}
	stream, err := c.PutFile(context.Background())
	// ...
}
```

## Extension interface

At present, some building interfaces such as SMS, encryption, email, etc in go sdk are not included in the `runtimeAPI` service, but instead use interface interfaces generated by proto. 
So when using this type of interface, you can directly call the grpc interface. 
Related interface documents refer to the proto file in the [spec directory](https://github.com/mosn/layotto/tree/main/spec/proto/extension/v1)

Taking the email interface as an example, the content is as follows:


> The interface is still under construction and there is no component implementation, so it cannot be truly called. The following is only a simple example 
 
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



## More Examples

For other SDK interfaces, please refer to the code examples in the [demo directory](https://github.com/mosn/layotto/tree/main/demo) and [quick start startup document](en/start/state/state.md)

Refer to the [configs example](https://github.com/mosn/layotto/tree/main/configs) for writing relevant configuration files and [Component Configuration Document](en/configuration/overview.md)