# 动态配置下发、组件热重载
## 1. 解决的问题
- 现在生产用户有一套定制的启动时初始化配置的方案：有一些配置配在 app 里，app 启动后、调 sidecar，让 sidecar 基于这些配置做初始化。方案不够通用，想做的更通用些

![image](https://user-images.githubusercontent.com/26001097/168947177-6a26397e-4648-47f0-a8df-e898285cd8f9.png)

- 支持"动态下发配置"。
    - 一种思路是让配置文件和镜像解耦，通过磁盘挂载进容器。比如 Dapr 的配置项放进 Configuration CRD, CRD 变更后，需要运维人员通过 k8s 滚动重启集群。
    - 另一种思路是[把 config_store 组件注入进别的组件](https://github.com/mosn/layotto/issues/500#issuecomment-1119390497) ，但有一些缺点:
        - 用户如果想用“动态配置下发”功能，没法当伸手党，没有社区现成的组件用，得开发自己的组件。
          最好是 runtime 层做一些通用功能，“赋能”所有组件，社区维护现成的组件、支持动态配置下发，方便用户当伸手党，开箱即用。
    - 另一种思路是像 envoy 一样，把配置分两类：bootstrap 配置（静态配置）、动态配置，前者放不放进镜像都行，后者支持配置动态下发、根据配置做热重载。

## 2. 产品设计
### User story
1. 用户在 apollo 页面改一下 Redis 的容灾切换配置，Redis 组件就能接收到新配置，把流量切到灾备集群
2. 已有生产用户可以把初始化流程迁移到新的模型，向下兼容。

### 编程界面
比如，现在 state.redis 的启动配置有下面这些（截图取自 [dapr 文档](https://docs.dapr.io/reference/components-reference/supported-state-stores/setup-redis/) )
![image](https://user-images.githubusercontent.com/26001097/168946975-9804d792-8851-463f-80ee-26231468f0aa.png)

现状是：redis 组件启动时，用这些配置kv做初始化；所有配置都是静态配置、只在启动时取一次，不监听后续配置变更。

但是我们可以改成：
- 这些 kv 可以动态下发
- layotto 监听这些 kv 的变更，一但有变化，**用最新的配置重新初始化组件**
- 如果组件觉得重新初始化太小题大做了，可以实现动态更新接口

优缺点分析：
- pros
    - runtime 层可以做一些通用功能，“赋能”所有组件；方便用户当伸手党，社区维护现成的组件、支持动态配置下发，用户开箱即用
- cons
    - 实现起来复杂。比如重新初始化期间，怎么保证流量无损？
    - 我不清楚这能不能满足用户生产需求，担心过早设计、过度设计

## 3. High-level design
![image](https://user-images.githubusercontent.com/26001097/168949648-3f440a84-45d3-45c1-89ef-79cb25d49713.png)

### 启动完成后，暴露 UpdateConfiguration API
Sidecar 启动还是用 json 文件，启动完成、readiness check 通过后，对外暴露一个新的 API，用于做配置热变更:

```protobuf
rpc UpdateConfiguration( RuntimeConfig) returns (UpdateResponse)
```

### Agent 负责和控制面交互、调用 UpdateConfiguration API
也就是说，Sidecar 只是开个接口、等别人推配置。而具体和控制面交互、订阅配置变更的事情可以封装 agent 来做，比如图上的 agent 2，负责订阅 apollo 的配置变更，有变更了就调 Sidecar 的接口，让 Sidecar 热更新。

对于已有的生产用户，可以像图上封装 agent 1, 监听 app 喂的配置、dump 配置、重启时加载配置，然后把配置推给 Sidecar。

再比如可以写个 File agent 问题，监听文件变化，有变化就读取新配置、通知 Sidecar 热重载。

agent 不一定要单独进程，在 main 里启动一个独立协程也行。

### 组件热重载
Sidecar 被调 UpdateConfiguration API 后，会:
1. 判断组件有没有实现"增量更新"接口:

```go
UpdateConfig(ctx context.Context, metadata map[string]string) (err error, needReload bool)
```

2. 如果组件有实现该接口，runtime 尝试让其增量更新
3. 如果增量更新失败，或者没实现该接口，则 runtime **根据全量配置重新初始化组件**
4. 新组件重新初始化完成后(通过 readiness check)，接管原组件的流量

## 4. 详细设计

### 4.1. gRPC API 设计

```protobuf
service Lifecycle {

  rpc ApplyConfiguration(DynamicConfiguration) returns (ApplyConfigurationResponse){}

}

message DynamicConfiguration{

  ComponentConfig component_config = 1;

}

message ApplyConfigurationResponse{
}
```

#### ComponentConfig 字段设计
##### a. 设计一个通用的更新接口

```protobuf
message ComponentConfig{

  // For example, `lock`, `state`
  string kind = 1;

  // The component name. For example,  `state_demo`
  string name = 2;

  map<string, string> metadata = 3;
}
```

~~用  google/protobuf/struct.proto  描述动态json 见 https://stackoverflow.com/questions/52966444/is-google-protobuf-struct-proto-the-best-way-to-send-dynamic-json-over-grpc~~

用 `map<string, string>` 传配置。

- 优点
  每次新加 API 或改配置结构时, 不用改每个语言的 sdk，让用户透传、sidecar 侧反序列化
  
- 缺点
  字段格式没有显示定义，不明确，不够结构化

##### b. 结构化定义每类配置

```protobuf
// Component configuration
message ComponentConfig{
  // For example, `lock`, `state`
  string kind = 1;
  // The component name. For example,  `state_demo`
  string name = 2;

  google.protobuf.Struct metadata = 3;

  oneof common_config {
    LockCommonConfiguration lock_config = 4;

    StateCommonConfiguration state_config = 5;

    // ....
  }
}
```

优缺点和上面相反

##### 结论
选择 A，减少 SDK 维护成本

#### Q: 是单独写一个 API 插件，还是放进已有的 API 插件里
单独写一个 API 插件

#### Q: 等人推配置 vs 主动拉配置 vs 推了之后再反拉
等人推配置

#### Q: API 接受全量配置还是增量配置
a. 增量，顺序问题由 stream 保证

```protobuf
service Lifecycle {

  rpc UpdateComponentConfiguration(stream ComponentConfig) returns (UpdateResponse){}

}
```

b. 全量

结论: b, 更简单。后面有需要的话可以再加一个通过 stream 做增量变更的接口。

### 4.2. 组件 API 设计

```go
type DynamicComponent interface {
    ApplyConfig(ctx context.Context, metadata map[string]string) (err error, needReload bool)
}
```

## 5. Future work
### pubsub 订阅关系下发

需要下发一些更结构化的配置数据

### 组件热重载
// TODO

- 重新初始化过程中，怎么保证流量无损
- 配置优先级：有一些配置是某个 app 定制的配置，有一些配置是所有 app 公用的通用配置，两者优先级是啥
- 配置事务读写，避免脏读
