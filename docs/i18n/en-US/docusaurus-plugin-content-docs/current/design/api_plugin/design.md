# API Plugin Design Documentation & Use Guide

The first half of this paper discusses what the API plugin solves and why it is designed, and the second half describes how the API plugin works.

If you use only API plugins, you can jump directly to [2.4.Use Guide](design/api_plugin/design?id=_24 - Use Guide)

## Needs analysis

### What issues to address

Resolve extension issues.The API, either Dapr or open source Layotto, is currently unable to fully meet production needs.

Looking at the development history of POSIX API and system call in the operating system domain, we can learn a lot to predict the future of Runtime.We can say that the Runtime API is unlikely to fully meet user needs in the future.Split OS Realm, even if there is a POSIX API, some scenes need to bypass standard
API, operate special hardware with special instructions.

Dapr extension is resolved through the Binding API, but this unstructured API has many problems (e.g. destroying portability, not supporting stream etc.)

### User scenarios and needs

For example, there are the following user scenes：

1. Companies have their own customized API needs because they are non-common and unsuited to open source Layotto/Dapr, so a company’s intermediate team wants to develop itself into sidecar.If a company's project import source Layotto or Dapr, there is no way to extend API development under the current architecture, only Fork
   to make an extension

![image](https://user-images.githubusercontent.com/26001097/131614836-60d797c8-b80b-4018-ad43-c2b874d35660.png)

User needs： in this case

- sdk sink;
- Multilingualism;
- Cloud deployment (just need intermediate teams themselves to develop components for cloud environments without community ready-made components)

2. Companies have new API requirements that fit into open source projects and thus demand the community, but the community has had difficulty in reaching consensus and has been struggling for months (e.g. https://github.com/dapr/dapr/issues/2988
   ).In such cases, companies may have operational pressures, cannot afford to wait for so long and want to relocate themselves before they move to new functions in the community.

User needs： in this case

- The user is autonomous about this feature and does not need to persuade the community (along with Chinese and English) to make this feature work for the world
- Rapid expansion, service operations

3. Users want to add fields to the Dapr API, first add fields to their Fork version, meet online needs, and then raise PR to the community.The community refused to add the field, PR was closed.用户很尴尬：这字段已经在线上使用了，怎么处理？

## High level design

### Hierarchical API

Refer to the OS Realm API for this year, we can design the Runtime API into multi-layer：

![img.png](https://gw.alipayobjects.com/mdn/rms_5891a1/afts/img/A*bWnHR7yhiF4AAAAAAAAAAAAAARQnAQ)

： for OS fields

- POSIX API
- Various Unix-like system of its own System Call (with portability, implementing the same interface with different hardware drivers)
- Special features provided by special hardware (no portable)

Based on this idea, we can design API plugins to support users extending their own private APIs.

![image](https://user-images.githubusercontent.com/26001097/131614802-c6f6a556-4e8b-4fee-b899-275a80e00eb6.png)

### Design objectives

1. Let open source users with customized development needs use Layotto instead of Fork by importing Layotto

2. Simple enough to develop api plugin

3. Config file is the same json, new api plugin does not need to add new configuration file

### Functional design

![image](https://user-images.githubusercontent.com/26001097/131614952-ccfc7d11-d376-41b0-b16c-2f17bfd2c9fc.png)

Layotto add a number of new extensions.

When using Layotto, a company user can maintain a project, import Layotto himself.Various extensions and components are stored in your project.If you are familiar with Java, this is similar to [Eureka]for Java communities (https://github.com/Netflix/eureka), you can
import Eureka and then expand.

When users want to add an API, they can develop a package in their own project (including their proto,pb files, their own API implementation), then call Layotto extension points in `main.go` and register their API in Layotto.

### Guide to use

How to add your own proto, add your own private API?

An example is [the shellowold package provided in the project](https://github.com/mosn/layotto/tree/main/cmd/layotto_multiple_api/helloworld) which implements custom API, \`\`Sayhello\`

Use this as an example to explain the steps to write an API plugin:

#### step 0. Define your own proto file and compile a pb

For example, users would like to add their own `Greeter` API, offering the `Sayhello` method, write a proto:
(this example is my paste from [grpc official exams](https://github.com/grpc/grpc-go/blob/master/examples/helloworld/helloworld/helloworld.proto)

```protobuf
syntax = "proto3";

option go_package = "google.golang.org/grpc/exames/hellotorld/hellowld";
option java_multiple_files = true;
option java_package = "io. rpc.examples.hellowld";
option java_outer_classname = "HelloWorldProto";

package helloord;

// The greening service definition.
Service Greeter {}
  // Sends a greening
  rpc Sayhello (HelloRequest) returns (HelloReply) {}
}

// The request message containing the user's name.
message HelloRequest {
  string name = 1;
}

// The response message containing the greies
message Hello Reply to the LO
  string message = 1;
}
```

Then compile it into the `.pb.go` file.

[Bellowold example pack provided by the project] (https://github.com/mosn/layotto/tree/main/cmd/layotto_multiple_api/helloworld) steal a lazy, direct import of grpc officially compiled .pb.go
file：

<img src="https://gw.alipayobjects.com/mdn/rms_5891a1/afts/img/A*9VnARJimj90AAAAAAAAAAAAAARQnAQ" width="40%" height="40%" alt="score" align="center" />

#### step 1. Write implementation for the API just defined

The protoc compiler will help you compile the interface `hellotorld.GreeterServer`, based on proto file, but the interface needs to be translated by itself.

For example, the `server` we wrote in the example implements `helloworl.GreeterServer` interface, with `Sayhello` method:

```go
//server is used to implement elllowld.GreeterServer.
type server struct LO
	appId strating
	// custom-components which implement the `HelloWorld` interface
	name2component map[string]component. World
	// LockStore components. They are not used in this demo, we put them here as a demo.
	name2LockStore map[string]block. ockStore
	pb.UnimplementedGreeterServer
}

// Sayhello implements hellowd. reeterServer.Sayhello
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) 6
	if _, ok := s.name2component[componentName]; !ok $6
		return &pb. elloReply{Message: "We don't want to talk with you!"}, nil
	}
	message, err := s.name2component[componentName]. ayHello(in.GetName()
	if err != nil {
		return nil, err
	}
	return &pb.HelloReply{Message: message}, nil
}
```

#### Step 2. Implement [`GrpcAPI` interface](https://github.com/mosn/layotto/blob/main/pkg/grpc/grpc_api.go), manage the life cycle of API plugins

Now that you have your own API implementation, you need to register it on Layotto next step.

> **Rememberly**：How to register the API on the original grpc server?
>
> Just write this line of code：
>
> pb.RegisterGreeterServer (s, &server{})

To register your API on Layotto,：

- Implementing [`GrpcAPI` interface](https://github.com/mosn/layotto/blob/main/pkg/grpc/grpc_api.go), implementing some lifecycle hooks

This GrpcAPI manages your API lifecycle and provides various life cycle hooks.There are currently Init and Register in the current lifecycle hook.

```go
// GrpcAPI is the interface of API plugin. It has lifecycle related methods
type GrpcAPI interface of the APIs
    // init this API before binding it to the grpc server.
    // For example, you can call the app to every their subscriptions.
    Init (conn *grpc. lientConn) error
    
    // Bind this API to the grpc server
    Register(s *grpc. erver, registeredServer mgrp.RegisteredServer) (mgrpc.RegisteredServer, error)
}
```

- Performs the corresponding constructor `NewGrpcAPI` to create your `GrpcAPI`.

```go
// NewGrpcAPI is the constructor of GpcAPI
type NewGrpcAPI func (applicationContext *ApplicationContext) GrpcAPI
```

`*ApplicationContext` is defined as：

```go
// ApplicationContext contains all you need to construct your GrpcAPI, such as all the components.
// For example, your `SuperState` GrpcAPI can hold the `StateStores` components and use them to implement your own `Super State API` logic.
type ApplicationContext struct {
    AppId                 string
    Hellos                map[string]hello.HelloService
    ConfigStores          map[string]configstores.Store
    Rpcs                  map[string]rpc.Invoker
    PubSubs               map[string]pubsub.PubSub
    StateStores           map[string]state.Store
    Files                 map[string]file.File
    LockStores            map[string]lock.LockStore
    Sequencers            map[string]sequencer.Store
    SendToOutputBindingFn func(name string, req *bindings.InvokeRequest) (*bindings.InvokeResponse, error)
    SecretStores          map[string]secretstores.SecretStore
    CustomComponent       map[string]map[string]custom.Component
}
```

##### What is the interpretation of：`CustomComponent`?

is "Custom components".

Component in Layotto is divided into prime mill：

- Preset Component

e.g. `pubsub` components, like `state` component

- Custom Component

Allows you to expand your own components, such as the `HelloWorld` component in the example below.

##### How does：configure a custom component?

See[自定义组件的配置文档](component_specs/custom/common.md)

##### View Example

Look at a specific example, in [helloowold exams](https://github.com/mosn/layotto/blob/main/cmd/layotto_multiple_api/helloworld/grpc_api.go), `*server` implements `Init`
and `Register`:

```go
func (s *server) Init (conn *rawGRPC.ClientConn) error {
	return nil
}

func (s *server) Register(grpcServer *rawGRPC.Server, registeredServer mgrpc. egisteredServer) (mgrpc.RegisteredServer, error) um
	pb.RegisterGreeterServer (GrpcServer, s)
	return registeredServer, nil
}
```

There is also a construction:

```go
func NewHelloWorldAPI(ac *grpc_api.ApplicationContext) grpc.GrpcAPI
	// 1. convert custom-components components
	name2component := make(map[string]component. elloWorld)
	if len(ac.CustomComponent) != 0 56
		/we only care about those components of type "helloowl"
		name2comp, ok := ac. customomComponent[kind]
		if ok & len (name2comp) > 0 F6
			for name, ::= range name2comp LO
				// convert them using type assortion
				comp, ok := v. component.HelloWorld)
				if !ok
					errMsg := fmt. printf("customom component %s does not implement HelloWorld interface", name)
					log.DefaultLogger. rorf(errMsg)
				}
				name2component[name] = comp
			}
		}
	}
	/ 2. Construct your API implementation
	return &server
		appId: ac. ppId,
		// Your API plugins in can store and use all the components.
		// For example, this demand set all the LockStore components here.
		name2LockStore: ac.LockStories,
		// Customers components of type "helloorl"
		name2component: name2component,
	}
}
```

##### Explain：these callbacks, constructions?

Look at this example, you might ask：these callbacks, constructions and calls?

The hook above is used to customize the start logic for the user extension.Layotto reverses the above life-cycle hooks and constructions during startup.Call order roughly：

`Layotto initialize all components` ---> Call `NewGrpcAPI` constructor ---> `GrpcAPI.Init` ---> ``Layotto create grpc server` --->``GrpcAPI.Register\`\`

Graph below：

<img src="https://gw.alipayobjects.com/mdn/rms_5891a1/afts/img/A*7_NyQL-FjigAAAAAAAAAAAAAARQnAQ" width="40%" height="40%" alt="score" align="center" />

#### step 3. Sign up your own API into Layotto

After achieving your private API
following the above steps, you can [register it in your main main in Layotto](https://github.com/mosn/layotto/blob/5234a80cdc97798162d03546eb8e0ee163c0ad60/cmd/layotto_multiple_api/main.go#L203):

```go

func NewRuntimeGrpcServer (data json.RawMessage, options..grpc.ServerOption) (mgrpc. egisteredServer, error) LO
	// ......
	
    /3. run
    server, err := rt. un(
        runtime.WithGrpcOptions(opts... ,
        // register your GrpcAPI here
        runtime. ithGrpcAPI(
            // default GrpcAPI
            default_api. ewGrpcAPI,
            // a demo to show how to register your own GrpcAPI
            helloord_api. ewHelloWorldAPI,
        ),
        // Hello
        runtime. ithHelloFactory (
            hello.NewHelloFactory ("hellotorld", hellold. ewHelloWorld),
        ),
    // ...
```

We recommend that users customize main functions and customize startup processes in their own projects.

Specifically, you can paste the main copy of Layotto into your project, modify it as necessary, and remove something that is not available (e.g. in a distribution lock component that is not etcd, you can remove it from your main name)

#### Step 4. Compile Run Layotto

Ready to start Layotto.

Example hellowd：

```shell
cd ${project_path}/cmd/layotto_multiple_api
go build -o layotto
# run it
./layotto start -c ../../configs/config_standalone.json
```

During Layotto launch, each registered API lifecycle method (Init, Register) will be callback to each registered API

Once startup, your API will offer grpc services externally
