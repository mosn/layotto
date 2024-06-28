# MOSN 子项目 Layotto：开启服务网格+应用运行时新篇章

>作者简介：
马振军，花名古今，在基础架构领域耕耘多年，对 Service Mesh 有深度实践经验，目前在蚂蚁集团中间件团队负责 MOSN、Layotto 等项目的开发工作。
Layotto 官方 GitHub 地址:[https://github.com/mosn/layotto](https://github.com/mosn/layotto)

点击链接即可观看现场视频：[https://www.bilibili.com/video/BV1hq4y1L7FY/](https://www.bilibili.com/video/BV1hq4y1L7FY/)

Service Mesh 在微服务领域已经非常流行，越来越多的公司开始在内部落地，蚂蚁从 Service Mesh 刚出现的时候开始，就一直在这个方向上大力投入，到目前为止，内部的 Mesh 方案已经覆盖数千个应用、数十万容器并且经过了多次大促考验，Service Mesh 带来的业务解耦，平滑升级等优势大大提高了中间件的迭代效率。

在大规模落地以后，我们又遇到了新的问题，本文主要对 Service Mesh 在蚂蚁内部落地情况进行回顾总结，并分享对 Service Mesh 落地后遇到的新问题的解决方案。

## 一、Service Mesh 回顾与总结

### A、Service Mesh 的初衷

>![](https://gw.alipayobjects.com/mdn/rms_1c90e8/afts/img/A*p8tGTbpLRegAAAAAAAAAAAAAARQnAQ)
在微服务架构下，基础架构团队一般会为应用提供一个封装了各种服务治理能力的 SDK，这种做法虽然保障了应用的正常运行，但缺点也非常明显，每次基础架构团队迭代一个新功能都需要业务方参与升级才能使用，尤其是 bugfix 版本，往往需要强推业务方升级，这里面的痛苦程度每一个基础架构团队成员都深有体会。

伴随着升级的困难，随之而来的就是应用使用的 SDK 版本差别非常大，生产环境同时跑着各种版本的 SDK，这种现象又会让新功能的迭代必须考虑各种兼容，就好像带着枷锁前进一般，这样随着不断迭代，会让代码维护非常困难，有些祖传逻辑更是一不小心就会掉坑里。

同时这种“重”SDK 的开发模式，导致异构语言的治理能力非常薄弱，如果想为各种编程语言都提供一个功能完整且能持续迭代的 SDK 其中的成本可想而知。

18 年的时候，Service Mesh 在国内持续火爆，这种架构理念旨在把服务治理能力跟业务解耦，让两者通过进程级别的通信方式进行交互。在这种架构模式下，服务治理能力从应用中剥离，运行在独立的进程中，迭代升级跟业务进程无关，这就可以让各种服务治理能力快速迭代，并且由于升级成本低，因此每个版本都可以全部升级，解决了历史包袱问题，同时 SDK 变“轻”直接降低了异构语言的治理门槛，再也不用为需要给各个语言开发相同服务治理能力的 SDK 头疼了。

### B、Service Mesh 落地现状

>![](https://gw.alipayobjects.com/mdn/rms_1c90e8/afts/img/A*rRG_TYlHMqYAAAAAAAAAAAAAARQnAQ)
蚂蚁很快意识到了 Service Mesh 的价值，全力投入到这个方向，用 Go 语言开发了 MOSN 这样可以对标 envoy 的优秀数据面，全权负责服务路由，负载均衡，熔断限流等能力的建设，大大加快了公司内部落地 Service Mesh 的进度。

现在 MOSN 在蚂蚁内部已经覆盖了数千个应用、数十万容器，新创建的应用默认接入 MOSN，形成闭环。而且在大家最关心的资源占用、性能损耗方面 MOSN 也交出了一份让人满意的答卷：
1.  RT 小于 0.2ms

2.  CPU 占用增加 0%~2%

3.  内存消耗增长小于 15M

由于 Service Mesh 降低了异构语言的服务治理门槛，NodeJS、C++等异构技术栈也在持续接入到 MOSN 中。

在看到 RPC 能力 Mesh 化带来的巨大收益之后，蚂蚁内部还把 MQ，Cache，Config 等中间件能力都进行了 Mesh 化改造，下沉到 MOSN，提高了中间件产品整体的迭代效率。

### C、新的挑战

1. 应用跟基础设施强绑定

>![](https://gw.alipayobjects.com/mdn/rms_1c90e8/afts/img/A*nKxcTKLp4EoAAAAAAAAAAAAAARQnAQ)
一个现代分布式应用，往往会同时依赖 RPC、Cache、MQ、Config 等各种分布式能力来完成业务逻辑的处理。

当初看到 RPC 下沉的红利以后，其他各种能力也都快速下沉。初期，大家都会以自己最熟悉的方式来开发，这就导致没有统一的规划管理，如上图所示，应用依赖了各种基础设施的 SDK，而每种 SDK 又以自己特有的方式跟 MOSN 进行交互，使用的往往都是由原生基础设施提供的私有协议，这直接导致了复杂的中间件能力虽然下沉，但应用本质上还是被绑定到了基础设施，比如想把缓存从 Redis 迁移到 Memcache 的话，仍旧需要业务方升级 SDK，这种问题在应用上云的大趋势下表现的更为突出，试想一下，如果一个应用要部署在云上，由于该应用依赖了各种基础设施，势必要先把整个基础设施搬到云上才能让应用顺利部署，这其中的成本可想而知。
因此如何让应用跟基础设施解绑，使其具备可移植能力，能够无感知跨平台部署是我们面临的第一个问题。

2. 异构语言接入成本高

>![](https://gw.alipayobjects.com/mdn/rms_1c90e8/afts/img/A*oIdQQZmgtyUAAAAAAAAAAAAAARQnAQ)
事实证明 Service Mesh 确实降低了异构语言的接入门槛，但在越来越多的基础能力下沉到 MOSN 以后，我们逐渐意识到为了让应用跟 MOSN 交互，各种 SDK 里都需要对通信协议，序列化协议进行开发，如果再加上需要对各种异构语言都提供相同的功能，那维护难度就会成倍上涨，

Service Mesh 让重 SDK 成为了历史，但对于现在各种编程语言百花齐放、各种应用又强依赖基础设施的场景来说，我们发现现有的 SDK 还不够薄，异构语言接入的门槛还不够低，如何进一步降低异构语言的接入门槛是我们面临的第二个问题。

## 二、Multi Runtime 理论概述

### A、什么是 Runtime?

>![](https://gw.alipayobjects.com/mdn/rms_1c90e8/afts/img/A*hQT-Spc5rI4AAAAAAAAAAAAAARQnAQ)
20 年初的时候，Bilgin lbryam 发表了一篇名为
Multi-Runtime Microservices Architecture
的文章，里面对微服务架构下一阶段的形态进行了讨论。

如上图所示，作者把分布式服务的需求进行了抽象，总共分为了四大类：
1.  生命周期（Lifecycle）
    主要指应用的编译、打包、部署等事情，在云原生的大趋势下基本被 docker、kubernetes 承包。

2.  网络（Networking）
    可靠的网络是微服务之间进行通信的基本保障，Service Mesh 正是在这方面做了尝试，目前 MOSN、envoy 等流行的数据面的稳定性、实用性都已经得到了充分验证。

3.  状态（State）
     分布式系统需要的服务编排，工作流，分布式单例，调度，幂等性，有状态的错误恢复，缓存等操作都可以统一归为底层的状态管理。

4.  绑定（Binding）
     在分布式系统中，不仅需要跟其他系统通信，还需要集成各种外部系统，因此对于协议转换，多种交互模型、错误恢复流程等功能也都有强依赖。

明确了需求以后，借鉴了 Service Mesh 的思路，作者对分布式服务的架构演进进行了如下总结：

>![](https://gw.alipayobjects.com/mdn/rms_1c90e8/afts/img/A*rwS2Q5yMp_sAAAAAAAAAAAAAARQnAQ)
第一阶段就是把各种基础设施能力从应用中剥离解耦，通通变成独立 sidecar 模型伴随着应用一起运行。

第二阶段是把各种 sidecar 提供的能力统一抽象成若干个 Runtime，这样应用从面向基础组件开发就演变成了面向各种分布式能力开发，彻底屏蔽掉了底层实现细节，而且由于是面向能力，除了调用提供各种能力的 API 之外，应用再也不需要依赖各种各样基础设施提供的 SDK 了。

作者的思路跟我们希望解决的问题一致，我们决定使用 Runtime 的理念来解决 Service Mesh 发展到现在所遇到的新问题。

### B、Service Mesh vs Runtime

>![](https://gw.alipayobjects.com/mdn/rms_1c90e8/afts/img/A*srPVSYTEHc4AAAAAAAAAAAAAARQnAQ)
为了让大家对 Runtime 有一个更加清晰的认识，上图针对 Service Mesh 跟 Runtime 两种理念的定位、交互方式、通信协议以及能力丰富度进行了总结，可以看到相比 Service Mesh 而言，Runtime 提供了语义明确、能力丰富的 API，可以让应用跟它的交互变得更加简单直接。

## 三、MOSN 子项目 Layotto

### A、dapr 调研

>![](https://gw.alipayobjects.com/mdn/rms_1c90e8/afts/img/A*Ab9HSYIK7CQAAAAAAAAAAAAAARQnAQ)
dapr 是社区中一款知名的 Runtime 实现产品，活跃度也比较高，因此我们首先调研了 dapr 的情况，发现 dapr 具有如下优势：
1.  提供了多种分布式能力，API 定义清晰，基本能满足一般的使用场景。

2.  针对各种能力都提供了不同的实现组件，基本涵盖了常用的中间件产品，用户可以根据需要自由选择。

当考虑如何在公司内部落地 dapr 时，我们提出了两种方案，如上图所示：

1.  替换：废弃掉现在的 MOSN，用 dapr 进行替换，这种方案存在两个问题：

a.  dapr 虽然提供了很多分布式能力，但目前并不具备 Service Mesh 包含的丰富的服务治理能力。

b.  MOSN 在公司内部已经大规模落地，并且经过了多次大促考验，直接用 dapr 来替换 MOSN 稳定性有待验证。

2.  共存：新增一个 dapr 容器，跟 MOSN 以两个 sidecar 的模式进行部署。这种方案同样存在两个问题：

a.  引入一个新的 sidecar，我们就需要考虑它配套的升级、监控、注入等等事情，运维成本飙升。

b.  多维护一个容器意味着多了一层挂掉的风险，这会降低现在的系统可用性。

同样的，如果你目前正在使用 envoy 作为数据面，也会面临上述问题。
因此我们希望把 Runtime 跟 Service Mesh 两者结合起来，通过一个完整的 sidecar 进行部署，在保证稳定性、运维成本不变的前提下，最大程度复用现有的各种 Mesh 能力。此外我们还希望这部分 Runtime 能力除了跟 MOSN 结合起来之外，未来也可以跟 envoy 结合起来，解决更多场景中的问题，Layotto 就是在这样的背景下诞生。

### B、Layotto 架构

>![](https://gw.alipayobjects.com/mdn/rms_1c90e8/afts/img/A*sdGoSYB_XFUAAAAAAAAAAAAAARQnAQ)
如上图所示，Layotto 是构建在 MOSN 之上，在下层对接了各种基础设施，向上层应用提供了统一的，具有各种各样分布式能力的标准 API。对于接入 Layotto 的应用来说，开发者不再需要关心底层各种组件的实现差异，只需要关注应用需要什么样的能力，然后调用对应能力的 API 即可，这样可以彻底跟底层基础设施解绑。

对应用来说，交互分为两块，一个是作为 gRPC Client 调用 Layotto 的标准 API，一个是作为 gRPC Server 来实现 Layotto 的回调，得利于gRPC 优秀的跨语言支持能力，应用不再需要关心通信、序列化等细节问题，进一步降低了异构技术栈的使用门槛。

除了面向应用，Layotto 也向运维平台提供了统一的接口，这些接口可以把应用跟 sidecar 的运行状态反馈给运维平台，方便 SRE 同学及时了解应用的运行状态并针对不同状态做出不同的举措，该功能考虑到跟 k8s 等已有的平台集成，因此我们提供了 HTTP 协议的访问方式。

除了 Layotto 本身设计以外，项目还涉及两块标准化建设，首先想要制定一套语义明确，适用场景广泛的 API 并不是一件容易的事情，为此我们跟阿里、 dapr 社区进行了合作，希望能够推进 Runtime API 标准化的建设，其次对于 dapr 社区已经实现的各种能力的 Components 来说，我们的原则是优先复用、其次开发，尽量不把精力浪费在已有的组件上面，重复造轮子。

最后 Layotto 目前虽然是构建在 MOSN 之上，未来我们希望 Layotto 可以跑在 envoy 上，这样只要应用接入了 Service Mesh，无论数据面使用的是 MOSN 还是 envoy，都可以在上面增加 Runtime能力。

### C、Layotto 的移植性

>![](https://gw.alipayobjects.com/mdn/rms_1c90e8/afts/img/A*2DrSQJ6GL8cAAAAAAAAAAAAAARQnAQ)
如上图所示，一旦完成 Runtime API 的标准化建设，接入 Layotto 的应用天然具备了可移植性，应用不需要任何改造就可以在私有云以及各种公有云上部署，并且由于使用的是标准 API，应用也可以无需任何改造就在 Layotto 跟 dapr 之间自由切换。

### D、名字含义

>![](https://gw.alipayobjects.com/mdn/rms_1c90e8/afts/img/A*CCckTZ_gRsMAAAAAAAAAAAAAARQnAQ)
从上面的架构图可以看出，Layotto 项目本身是希望屏蔽基础设施的实现细节，向上层应用统一提供各种分布式能力，这种做法就好像是在应用跟基础设施之间加了一层抽象，因此我们借鉴了 OSI 对网络定义七层模型的思路，希望 Layotto 可以作为第八层对应用提供服务，otto 是意大利语中8的意思，Layer otto 就是第八层的意思，简化了一下变成了 Layotto，同时项目代号 L8，也是第八层的意思，这个代号也是设计我们项目 LOGO 时灵感的来源。

介绍完项目的整体情况，下面对其中四个主要功能的实现细节进行说明。

### E、配置原语

>![](https://gw.alipayobjects.com/mdn/rms_1c90e8/afts/img/A*mfkRQZH3oNwAAAAAAAAAAAAAARQnAQ)
首先是分布式系统中经常使用的配置功能，应用一般使用配置中心来做开关或者动态调整应用的运行状态。Layotto 中配置模块的实现包括两部分，一个是对如何定义配置这种能力的 API 的思考，一个是具体的实现，下面逐个来看。

想要定义一个能满足大部分实际生产诉求的配置 API 并不是一件容易的事，dapr 目前也缺失这个能力，因此我们跟阿里以及 dapr 社区一起合作，为如何定义一版合理的配置 API 进行了激烈讨论。

目前讨论结果还没有最终确定，因此 Layotto 是基于我们提给社区的第一版草案进行实现，下面对我们的草案进行简要说明。

我们先定义了一般配置所需的基本元素：

1.  appId：表示配置属于哪个应用

2.  key：配置的 key

3.  content：配置的值

4.  group：配置所属的分组，如果一个 appId 下面的配置过多，我们可以给这些配置进行分组归类，便于维护。

此外我们追加了两种高级特性，用来适配更加复杂的配置使用场景：

1.  label，用于给配置打标签，比如该配置属于哪个环境，在进行配置查询的时候，我们会使用 label + key 来查询配置。

2.  tags，用户给配置追加的一些附加信息，如描述信息、创建者信息，最后修改时间等等，方便配置的管理，审计等。

对于上述定义的配置 API 的具体实现，目前支持查询、订阅、删除、创建、修改五种操作，其中订阅配置变更后的推送使用的是 gRPC 的 stream 特性，而底层实现这些配置能力的组件，我们选择了国内流行的 apollo，后面也会根据需求增加其他实现。

### F、Pub/Sub 原语

>![](https://gw.alipayobjects.com/mdn/rms_1c90e8/afts/img/A*YJs-R6WFhkgAAAAAAAAAAAAAARQnAQ)
对于 Pub/Sub 能力的支持，我们调研了 dapr 现在的实现，发现基本上已经可以满足我们的需求，因此我们直接复用了 dapr 的 API 以及 components，只是在 Layotto 里面做了适配，这为我们节省了大量的重复劳动，我们希望跟 dapr 社区保持一种合作共建的思路，而不是重复造轮子。

其中 Pub 功能是 App 调用 Layotto 提供的 PublishEvent 接口，而 Sub 功能则是应用通过 gRPC Server 的形式实现了 ListTopicSubscriptions 跟 OnTopicEvent 两个接口，一个用来告诉 Layotto 应用需要订阅哪些 topic，一个用于接收 topic 变化时 Layotto 的回调事件。

dapr 对于 Pub/Sub 的定义基本满足我们的需求，但在某些场景下仍有不足，dapr 采用了 CloudEvent 标准，因此 Pub 接口没有返回值，这无法满足我们生产场景中要求 Pub 消息以后服务端返回对应的 messageID 的需求，这一点我们已经把需求提交给了 dapr 社区，还在等待反馈，考虑到社区异步协作的机制，我们可能会先社区一步增加返回结果，然后再跟社区探讨一种更好的兼容方案。

### G、RPC 原语

>![](https://gw.alipayobjects.com/mdn/rms_1c90e8/afts/img/A*i-JnSaeZbJ4AAAAAAAAAAAAAARQnAQ)
RPC 的能力大家不会陌生，这可能是微服务架构下最最基础的需求，对于 RPC 接口的定义，我们同样参考了 dapr 社区的定义，发现完全可以满足我们的需求，因此接口定义就直接复用 dapr 的，但目前 dapr 提供的 RPC 实现方案还比较薄弱，而 MOSN 经过多年迭代，能力已经非常成熟完善，因此我们大胆把 Runtime 跟 Service Mesh 两种思路结合在一起，把 MOSN 本身作为我们实现 RPC 能力的一个 Component，这样 Layotto 在收到 RPC 请求以后交给 MOSN 进行实际数据传输，这种方案可以通过 istio 动态改变路由规则，降级限流等等设置，相当于直接复用了 Service Mesh 的各种能力，这也说明 Runtime 不是要推翻 Service Mesh，而是要在此基础上继续向前迈一步。

具体实现细节上，为了更好的跟 MOSN 融合，我们在 RPC 的实现上面加了一层 Channel，默认支持dubbo，bolt，http 三种常见的 RPC 协议，如果仍然不能满足用户场景，我们还追加了 Before/After 两种 Filter，可以让用户做自定义扩展，实现协议转换等需求。

### H、Actuator 原语

>![](https://gw.alipayobjects.com/mdn/rms_1c90e8/afts/img/A*E_Q-T4d_bm4AAAAAAAAAAAAAARQnAQ)
在实际生产环境中，除了应用所需要的各种分布式能力以外，PaaS 等运维平台往往需要了解应用的运行状态，基于这种需求，我们抽象了一套 Actuator 接口，目前 dapr 还没有提供这方面的能力，因此我们根据内部的需求场景进行了设计，旨在把应用在启动期、运行期等阶段各种各样的信息暴露出去，方便 PaaS 了解应用的运行情况。

Layotto 把暴露信息分为两大类：

1.  Health：该模块判断应用当前运行状态是否健康，比如某个强依赖的组件如果初始化失败就需要表示为非健康状态，而对于健康检查的类型我们参考了 k8s，分为：

a.  Readiness：表示应用启动完成，可以开始处理请求。

b.  Liveness：表示应用存活状态，如果不存活则需要切流等。

2.  Info：该模块预期会暴露应用的一些依赖信息出去，如应用依赖的服务，订阅的配置等等，用于排查问题。

Health 对外暴露的健康状态分为以下三种：

1.  INIT：表示应用还在启动中，如果应用发布过程中返回该值，这个时候 PaaS 平台应该继续等待应用完成启动。

2.  UP：表示应用启动正常，如果应用发布过程中返回该值，意味着 PasS 平台可以开始放入流量。

3.  DOWN：表示应用启动失败，如果应用发布过程中返回该值，意味着 PaaS 需要停止发布并通知应用 owner。

到这里关于 Layotto 目前在 Runtime 方向上的探索基本讲完了，我们通过定义明确语义的 API，使用 gRPC 这种标准的交互协议解决了目前面临的基础设施强绑定、异构语言接入成本高两大问题。随着未来 API 标准化的建设，一方面可以让接入 Layotto 的应用无感知的在各种私有云、公有云上面部署，另一方面也能让应用在 Layotto，dapr 之间自由切换，提高研发效率。

目前 Serverless 领域也是百花齐放，没有一种统一的解决方案，因此 Layotto 除了在上述 Runtime 方向上的投入以外，还在 Serverless 方向上也进行了一些尝试，下面就尝试方案进行介绍。

## 四、WebAssembly 的探索

### A、WebAssembly 简介

>![](https://gw.alipayobjects.com/mdn/rms_1c90e8/afts/img/A*-ACRSpqbuJ0AAAAAAAAAAAAAARQnAQ)
WebAssembly，简称 WASM，是一个二进制指令集，最初是跑在浏览器上来解决 JavaScript 的性能问题，但由于它良好的安全性，隔离性以及语言无关性等优秀特性，很快人们便开始让它跑在浏览器之外的地方，随着 WASI 定义的出现，只需要一个 WASM 运行时，就可以让 WASM 文件随处执行。

既然 WebAssembly 可以在浏览器以外的地方运行，那么我们是否能把它用在 Serverless 领域？目前已经有人在这方面做了一些尝试，不过如果这种方案真的想落地的话，首先要考虑的就是如何解决运行中的 WebAssembly 对各种基础设施的依赖问题。

### B、WebAssembly 落地原理

目前 MOSN 通过集成 WASM Runtime 的方式让 WASM 跑在 MOSN 上面，以此来满足对 MOSN 做自定义扩展的需求。同时，Layotto 也是构建在 MOSN 之上，因此我们考虑把二者结合在一起，实现方案如下图所示：

>![](https://gw.alipayobjects.com/mdn/rms_1c90e8/afts/img/A*U7UDRYyBOvIAAAAAAAAAAAAAARQnAQ)
开发者可以使用 Go/C++/Rust 等各种各样自己喜欢的语言来开发应用代码，然后把它们编译成 WASM 文件跑在 MOSN 上面，当 WASM 形态的应用在处理请求的过程中需要依赖各种分布式能力时就可以通过本地函数调用的方式调用 Layotto 提供的标准 API，这样直接解决了 WASM 形态应用的依赖问题。

目前 Layotto 提供了 Go 跟 Rust 版 WASM 的实现，虽然只支持 demo 级功能，但已经足够让我们看到这种方案的潜在价值。

此外，WASM 社区目前还处于初期阶段，有很多地方需要完善，我们也给社区提交了一些 PR共同建设，为 WASM 技术的落地添砖加瓦。

### C、WebAssembly 落地展望

>![](https://gw.alipayobjects.com/mdn/rms_1c90e8/afts/img/A*NzwKRY2GZPcAAAAAAAAAAAAAARQnAQ)
虽然现在 Layotto 中对 WASM 的使用还处于试验阶段，但我们希望它最终可以成为 Serverless 的一种实现形态，如上图所示，应用通过各种编程语言开发，然后统一编译成 WASM 文件，最后跑在 Layotto+MOSN 上面，而对于应用的运维管理统一由 k8s、docker、prometheus 等产品负责。

## 五、社区规划

最后来看下 Layotto 在社区的做的一些事情。

### A、Layotto vs Dapr

>![](https://gw.alipayobjects.com/mdn/rms_1c90e8/afts/img/A*OpQTRqoMpK0AAAAAAAAAAAAAARQnAQ)
上图列出了 Layotto 跟 dapr 现有的能力对比，在 Layotto 的开发过程中，我们借鉴 dapr 的思路，始终以优先复用、其次开发为原则，旨在达成共建的目标，而对于正在建设或者未来要建设的能力来说，我们计划优先在 Layotto 上落地，然后再提给社区，合并到标准 API，鉴于社区异步协作的机制，沟通成本较高，因此短期内可能 Layotto 的 API 会先于社区，但长期来看一定会统一。

### B、API 共建计划

>![](https://gw.alipayobjects.com/mdn/rms_1c90e8/afts/img/A*GAe8QqZ03eoAAAAAAAAAAAAAARQnAQ)
关于如何定义一套标准的 API 以及如何让 Layotto 可以跑在 envoy 上等等事项，我们已经在各个社区进行了深入讨论，并且以后也还会继续推进。

### C、Roadmap

>![](https://gw.alipayobjects.com/mdn/rms_1c90e8/afts/img/A*iUV3Q7S3VLEAAAAAAAAAAAAAARQnAQ)
Layotto 在目前主要支持 RPC、Config、Pub/Sub、Actuator 四大功能，预计在九月会把精力投入到分布式锁、State、可观测性上面，十二月份会支持 Layotto 插件化，也就是让它可以跑在 envoy 上，同时希望对 WebAssembly 的探索会有进一步的产出。

### D、正式开源

>![](https://gw.alipayobjects.com/mdn/rms_1c90e8/afts/img/A*S6mdTqAapLQAAAAAAAAAAAAAARQnAQ)
前面详细介绍了 Layotto 项目，最重要的还是该项目今天作为 MOSN 的子项目正式开源，我们提供了详细的文档以及 demo 示例方便大家快速上手体验。

对于 API 标准化的建设是一件需要长期推动的事情，同时标准化意味着不是满足一两种场景，而是尽可能的适配大多数使用场景，为此我们希望更多的人可以参与到 Layotto 项目中，描述你的使用场景，讨论 API 的定义方案，一起提交给社区，最终达成 Write once, Run anywhere 的终极目标！
