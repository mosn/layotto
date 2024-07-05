# 新增API时的开发规范
感谢您对Layotto的支持！

本文档描述了如何设计并实现新的Layotto API。Layotto用Go语言编写，如果您对Go语言不熟悉可以看下[Go教程](https://tour.golang.org/welcome/1) 。

在开发新API的时候，可以参考已有的其他API的代码、文档，这会让开发简单很多。

**Q: 为啥要制定该规范？**

A: 目前缺少使用文档，用户不好用，例如：

![img_1.png](../../img/development/api/img_1.png)

![img_2.png](../../img/development/api/img_2.png)

代码缺少注释，感兴趣的贡献者看不懂，例如 https://github.com/mosn/layotto/issues/112

旧文档和注释补起来太慢，希望今后开发的新功能有这些。

**Q: 遵循规范太麻烦了，会不会让想贡献代码的同学望而却步？**

A: **本规范只限制“新增Layotto API的pr需要有哪些东西”（比如新设计一个分布式自增id API)** ，其他pr比如新开发一个组件、新开发一个sdk都不需要遵循本规范，没这么复杂，足够自由

## 太长不看
开发前先提案，提案要详细

开发时要写4个给用户看的文档
- Quick start
- 使用文档
- API通用配置
- 组件配置

不用写设计文档，但是proto API和组件API要写详细注释，注释 as doc

新增API的pr要两个人code review，后续有机器人了可以一个人cr；其他pr随意

## 一、向社区发布API提案，经过充分讨论
### 1.1. 发布详细的提案
#### 1.1.1. 为什么提案要详细
如果提案粒度太粗，其他人评审时可能没啥好评的，发现不了问题；

评审的目的是集思广益，大家一起帮忙分析当前的设计存在的不足，尽早暴露问题，免得以后返工。

#### 1.1.2. 提案的内容
提案需要包含以下内容：

- 需求分析
  - 为什么要做这个API
  - 定义需求的边界，哪些feature支持，哪些不支持  
- 市面上产品调研  
- grpc/http API设计
- 组件API设计  
- 解释你的设计

一个优秀的提案示例：https://github.com/dapr/dapr/issues/2988

### 1.2. 提案评审
简单的API发出来后大家文字讨论即可；

重要或复杂的API设计可以组织社区会议进行评审。


## 二、开发
### 2.1. 代码规范

### 2.2. 测试规范
- 有单元测试
- 有client demo，可以拿来做演示、当集成测试

### 2.3. 文档规范
原则：需要写给用户看的文档；至于给开发者看的设计文档，因为时间长了后可能过期、和代码不一致，可以不写，通过贴proposal issue的链接、在代码里写注释来解释设计。

#### 2.3.1. Quick start
需要有：

- 这API是干嘛的
- 这个quickstart是干嘛的，想实现啥效果，最好有个图解释下
- 操作步骤

正例：[Dapr pub-sub quickstart](https://github.com/dapr/quickstarts/tree/v1.0.0/pub-sub) 在操作之前贴图解释下要做什么事情

![img.png](../../img/development/api/img.png)

反例：文档只写了操作步骤1234，用户看不懂操作这些想干啥

#### 2.3.2. 使用文档
文档路径在"用户手册--接口文档"下，例如 State API的见 https://mosn.io/layotto/#/zh/api_reference/state/reference

>调研发现Dapr的使用文档较多，比如光State API就有:
> 
> https://docs.dapr.io/developing-applications/building-blocks/state-management/  
> https://docs.dapr.io/reference/api/state_api/ 
> 
> https://docs.dapr.io/operations/components/setup-state-store/
> 
> https://docs.dapr.io/reference/components-reference/supported-state-stores/
> 
> 我们处于项目早期，可以轻一些

需要有：
##### what.这个API是啥，解决啥问题
##### when.什么场景适合用这个API
##### how.怎么用这个API
- 接口列表。例如：
  
![img_4.png](../../img/development/api/img_4.png)
  
列出来有哪些接口，一方面省的用户自己去翻proto、不知道哪些是相关API,一方面避免用户产生"这项目连接口文档都没有？！"的反感
- 关于接口的出入参：拿proto注释当接口文档  
考虑到接口文档用中英文写要写两份、时间长了还有可能和代码不一致，因此建议不写接口文档，直接把proto注释写的足够详细、当接口文档。例如：

```protobuf
// GetStateRequest is the message to get key-value states from specific state store.
message GetStateRequest {
  // Required. The name of state store.
  string store_name = 1;

  // Required. The key of the desired state
  string key = 2;

  // (optional) read consistency mode
  StateOptions.StateConsistency consistency = 3;

  // (optional) The metadata which will be sent to state store components.
  map<string, string> metadata = 4;
}
   
// StateOptions configures concurrency and consistency for state operations
message StateOptions {
  // Enum describing the supported concurrency for state.
  // The API server uses Optimized Concurrency Control (OCC) with ETags.
  // When an ETag is associated with an save or delete request, the store shall allow the update only if the attached ETag matches with the latest ETag in the database.
  // But when ETag is missing in the write requests, the state store shall handle the requests in the specified strategy(e.g. a last-write-wins fashion).
  enum StateConcurrency {
    CONCURRENCY_UNSPECIFIED = 0;
    // First write wins
    CONCURRENCY_FIRST_WRITE = 1;
    // Last write wins
    CONCURRENCY_LAST_WRITE = 2;
  }
  // Enum describing the supported consistency for state.
  enum StateConsistency {
    CONSISTENCY_UNSPECIFIED = 0;
    //  The API server assumes data stores are eventually consistent by default.A state store should:
    //
    // - For read requests, the state store can return data from any of the replicas
    // - For write request, the state store should asynchronously replicate updates to configured quorum after acknowledging the update request.
    CONSISTENCY_EVENTUAL = 1;

    // When a strong consistency hint is attached, a state store should:
    //
    // - For read requests, the state store should return the most up-to-date data consistently across replicas.
    // - For write/delete requests, the state store should synchronisely replicate updated data to configured quorum before completing the write request.
    CONSISTENCY_STRONG = 2;
  }

  StateConcurrency concurrency = 1;
  StateConsistency consistency = 2;
}
```      

这就要求proto注释里写清楚：
- 是必传参数还是可选参数；
- 解释这个字段啥含义；光解释字面意思是不够的，要解释背后的使用机制，比如上面的consistency和concurrency要解释用户传了某个选项后，服务器能提供什么样的保证
  
（consistency和concurrency上面的注释其实是我把Dapr文档上的描述精简后粘过来的，省了写双语文档）
- 注释讲不清楚的，在文档上解释

##### why.为什么这么设计
有设计文档的话贴个文档链接，没文档的话贴个proposal issue链接

#### 2.3.3. 介绍API通用配置的文档
例如https://mosn.io/layotto/#/zh/component_specs/state/common

- 配置文件结构
- 解释这个API的通用配置，比如keyPrefix

#### 2.3.4. 介绍组件配置的文档
例如https://mosn.io/layotto/#/zh/component_specs/state/redis

- 这个组件的配置项说明
- 想启动这个组件跑demo的话，怎么启动

### 2.4. 注释规范
#### proto注释 as doc

见上

#### 组件API 注释 as doc

如果不写双语设计文档，那么组件API的注释要承担设计文档的作用（向其他开发者解释）。可以贴下proposal issue的链接

判断写得好不好的标准是"发出去后，社区爱好者想贡献组件的话，能否不当面提问、自己看项目就能上手开发组件"

如果觉得注释解释不清楚，就写个设计文档，或者补充下proposal issue、写的更详细些吧

#### 其他注意事项

确保没有中文注释；

不用写无意义注释（把方法名复述一遍），比如：

```protobuf
	//StopSubscribe stop subs
	StopSubscribe()
```

## 三、提交pull request
### 3.1. 不符合开发规范的pr不可以合并进主干

### 3.2. cr人数
新增API的code review需要两个人review，后续有机器人自动检查后改成1个人review。

其他pull request的cr人数随意，不做约束。