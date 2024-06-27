# 蚂蚁云原生应用运行时的探索和实践 - ArchSummit 上海

>Mesh 模式的引入是实现应用云原生的关键路径，蚂蚁集团已在内部实现大规模落地。随着 Message、DB、Cache Mesh 等更多的中间件能力的下沉，从 Mesh 演进而来的应用运行时将是中间件技术的未来形态。应用运行时旨在帮助开发人员快速的构建云原生应用，帮助应用和基础设施进一步解耦，而应用运行时最核心是 API 标准，期望社区一起共建。

>![](https://gw.alipayobjects.com/mdn/rms_1c90e8/afts/img/A*nergRo8-RI0AAAAAAAAAAAAAARQnAQ)
## 蚂蚁集团 Mesh 化介绍

蚂蚁是一家技术和创新驱动的公司，从最早淘宝里的一个支付应用，到现在服务
全球十二亿用户的大型公司，蚂蚁的技术架构演进大概会分为如下几个阶段：

2006 之前，最早的支付宝就是一个集中式的单体应用，不同的业务做了模块化的开发。

2007 年的时候，随着更多场景支付的推广，开始做了一下应用、数据的拆分，做了 SOA 化的一些改造。

2010 年之后，推出了快捷支付，移动支付，支撑双十一，还有余额宝等现象级产品，用户数到了亿这个级别，蚂蚁的应用数也数量级的增长，蚂蚁自研了很多全套的微服务中间件去支撑蚂蚁的业务；

2014 年，像借呗花呗、线下支付、更多场景更多业务形态的出现，对蚂蚁的可用性和稳定性提出更高的要求，蚂蚁对微服务中间件进行了 LDC 单元化的支持，支撑业务的异地多活，以及为了支撑双十一超大流量的混合云上的弹性扩缩容。

2018 年，蚂蚁的业务不仅仅是数字金融，还有数字生活、国际化等一些新战略的出现，促使我们要有更加高效的技术架构能让业务跑得更快更稳，所以蚂蚁结合业界比较流行的云原生的理念，在内部进行了 Service Mesh、Serverless、可信原生方向的一些落地。

>![](https://gw.alipayobjects.com/mdn/rms_1c90e8/afts/img/A*KCSVTZWSf8wAAAAAAAAAAAAAARQnAQ)

可以看到蚂蚁的技术架构也是跟随公司的业务创新不断演进的，前面的从集中式到 SOA 再到微服务的过程，相信搞过微服务的同学都深有体会，而从微服务到云原生的实践是蚂蚁近几年自己探索出来的。

## 为什么要引入 Service Mesh

蚂蚁既然有一套完整的微服务治理中间件，那为什么还需要引入 Service Mesh 呢？

>![](https://gw.alipayobjects.com/mdn/rms_1c90e8/afts/img/A*Sq7oR6eO2QAAAAAAAAAAAAAAARQnAQ)

拿蚂蚁自研的服务框架 SOFARPC 为例，它是一个功能强大的 SDK，包含了服务发现、路由、熔断限流等一系列能力。在一个基本的 SOFA(Java) 应用里，业务代码集成了 SOFARPC 的 SDK，两者在一个进程里运行。在蚂蚁的大规模落地微服务之后，我们就面临了如下的一些问题：

**升级成本高**：SDK 是需要业务代码引入的，每次的升级都需要应用修改代码进行发布。由于应用规模较大，在一些大的技术变更或者安全问题修复的时候。每次需要数千个应用一起升级，费时费力。
**版本碎片化**：由于升级成本高，SDK 版本碎片化严重，这就导致我们写代码的时候需要兼容历史逻辑，整体技术演进困难。
**跨语言无法治理**：蚂蚁的中后台在线应用大多使用 Java 作为技术栈，但是在前台、AI、大数据等领域有很多的跨语言应用，例如 C++/Python/Golang 等等，由于没有对应语言的 SDK，他们的服务治理能力其实是缺失的。

我们注意到云原生里有 Service Mesh 一些理念开始出现，所以开始往这个方向探索。在 Service Mesh 的理念里，有两个概念，一个是 Control Plane 控制平面，一个是 Data Plane 数据平面。控制面这里暂时不展开，其中数据平面的核心思想就是解耦，将一些业务无需关系的复杂逻辑（如 RPC 调用里的服务发现、服务路由、熔断限流、安全）抽象到一个独立进程里去。只要保持业务和独立进程的通信协议不变，这些能力的演进可以跟随这个独立的进程自主升级，整个 Mesh 就可以做到统一演进。而我们的跨语言应用，只要流量是经过我们的 Data Plane 的，都可以享受到刚才提到的各种服务治理相关的能力，应用对底层的基础设施能力是透明的，真正的云原生的。

## 蚂蚁 Mesh 落地过程

所以从 2017 年底开始，蚂蚁就开始探索 Service Mesh 的技术方向，并提出了 基础设施统一，业务无感升级 的愿景。主要的里程碑就是：

2017 年底开始技术预研 Service Mesh 技术，并确定为未来发展方向；

2018 年初开始用 Golang 自研 Sidecar MOSN 并开源，主要支持 RPC 在双十一小范围试点；

2019 年 618，新增 Message Mesh 和 DB Mesh 的形态，覆盖若干核心链路，支撑 618 大促；

2019 年双十一，覆盖了所有大促核心链路几百个应用，支撑当时的双十一大促；

2020 年双十一，全站超过 80% 的在线应用接入了 Mesh 化，整套 Mesh 体系也具备了 2 个月从能力开发到全站升级完成的能力。

## 蚂蚁 Mesh 落地架构

目前 Mesh 化在蚂蚁落地规模是应用约数千个，容器数十万的级别，这个规模的落地，在业界是数一数二的，根本就没有前人的路可以学习，所以蚂蚁在落地过程中，也建设一套完整的研发运维体系去支撑蚂蚁的 Mesh 化。

>![](https://gw.alipayobjects.com/mdn/rms_1c90e8/afts/img/A*eAlMT7SMTpMAAAAAAAAAAAAAARQnAQ)

蚂蚁 Mesh 架构大概如图所示，底下是我们的控制平面，里面部署了服务治理中心、PaaS、监控中心等平台的服务端，都是现有的一些产品。还有就是我们的运维体系，包括研发平台和 PaaS 平台。那中间是我们的主角数据平面 MOSN，里面管理了 RPC、消息、MVC、任务四种流量，还有健康检查、监控、配置、安全、技术风险都下沉的基础能力，同时 MOSN 也屏蔽了业务和基础平台的一些交互。DBMesh 在蚂蚁是一个独立的产品，图里就没画出来。然后最上层是我们的一些应用，目前支持 Java、Nodejs 等多种语言的接入。
对应用来说，Mesh 虽然能做到基础设施解耦，但是接入还是需要一次额外的升级成本，所以为了推进应用的接入，蚂蚁做了整个研发运维流程的打通，包括在现有框架上做最简化的接入，通过分批推进把控风险和进度，让新应用默认接入 Mesh 化等一些事情。

同时随着下沉能力的越来越多，各个能力之前也面临了研发协作的一些问题，甚至互相影响性能和稳定性的问题，所以对于 Mesh 自身的研发效能，我们也做了一下模块化隔离、新能力动态插拔、自动回归等改进，目前一个下沉能力从开发到全站推广完成可以在 2 个月内完成。

## 云原生应用运行时上的探索

**大规模落地后的新问题与思考**

蚂蚁 Mesh 大规模落地之后，目前我们遇到了一些新的问题：
跨语言 SDK 的维护成本高：拿 RPC 举例，大部分逻辑已经下沉到了 MOSN 里，但是还有一部分通信编解码协议的逻辑是在 Java 的一个轻量级 SDK 里的，这个 SDK 还是有一定的维护成本的，有多少个语言就有多少个轻量级 SDK，一个团队不可能有精通所有语言的研发，所以这个轻量级 SDK 的代码质量就是一个问题。

业务兼容不同环境的新场景：蚂蚁的一部分应用是既部署在蚂蚁内部，也对外输出到金融机构的。当它们部署到蚂蚁时，对接的是蚂蚁的控制面，当对接到银行的时候，对接的是银行已有的控制面。目前大多数应用的做法是自己在代码里封装一层，遇到不支持的组件就临时支持对接一下。

从 Service Mesh 到 Multi-Mesh：蚂蚁最早的场景是 Service Mesh，MOSN 通过网络连接代理的方式进行了流量拦截，其它的中间件都是通过原始的 SDK 与服务端进行交互。而现在的 MOSN 已经不仅仅是 Service Mesh 了，而是 Multi-Mesh，因为除了 RPC，我们还支持了更多中间件的 Mesh 化落地，包括消息、配置、缓存的等等。可以看到每个下沉的中间件，在应用侧几乎都有一个对应的轻量级 SDK 存在，这个在结合刚才的第一问题，就发现有非常多的轻量级 SDK 需要维护。为了保持功能不互相影响，每个功能它们开启不同的端口，通过不同的协议去和 MOSN 进行调用。例如 RPC 用的 RPC 协议，消息用的 MQ 协议，缓存用的 Redis 协议。然后现在的 MOSN 其实也不仅仅是面向流量了，例如配置就是暴露了一下 API 给业务代码去使用。

>![](https://gw.alipayobjects.com/mdn/rms_1c90e8/afts/img/A*80o8SYwyHJoAAAAAAAAAAAAAARQnAQ)

为了解决刚才的问题和场景，我们就在思考如下的几个点：

1.不同中间件、不同语言的 SDK 能否风格统一？

2.各个下沉能力的交互协议能否统一？

3.我们的中间件下沉是面向组件还是面向能力？

4.底层的实现是否可以替换？

>![](https://gw.alipayobjects.com/mdn/rms_1c90e8/afts/img/A*hsZBQJg0VnoAAAAAAAAAAAAAARQnAQ)

## 蚂蚁云原生应用运行时架构

从去年的 3 月份开始，经过内部的多轮讨论，以及对业界一些新理念的调研，我们提出了一个“云原生应用运行时”（下称运行时）的概念。顾名思义，我们希望这个运行时能够包含应用所关心的所有分布式能力，帮助开发人员快速的构建云原生应用，帮助应用和基础设施进一步解耦！

>![](https://gw.alipayobjects.com/mdn/rms_1c90e8/afts/img/A*iqQoTYAma4YAAAAAAAAAAAAAARQnAQ)

云原生应用运行时设计里核心的几个点如下：

**第一**，由于有了 MOSN 规模化落地的经验和配套的运维体系，我们决定基于 MOSN 内核去开发我们的云原生应用运行时。

**第二**，面向能力，而不是面向组件，统一定义出这个运行时的 API 能力。

**第三**，业务代码和 Runtime API 之间的交互采用统一的 gRPC 协议，这样的话，业务端侧可以直接通过 proto 文件去反向生成一个客户端，直接进行调用。

**第四**，能力后面对应的组件实现是可以替换的，例如注册服务的提供者可以是 SOFARegistry，也可以是 Nacos 或者 Zookeeper。


**运行时能力抽象**

>![](https://gw.alipayobjects.com/mdn/rms_1c90e8/afts/img/A*hWIVR6ccduYAAAAAAAAAAAAAARQnAQ)

为了抽象出云原生应用最需要的一些能力，我们先定了几个原则：

1.关注分布式应用所需的 API 和场景而不是组件；
2.API 符合直觉，开箱即用，约定优于配置；
3.API 不绑定实现，实现差异化使用扩展字段。

有了原则之后，我们就抽象出了三组 API，分别是应用调用运行时的 mosn.proto，运行时调用应用的 appcallback.proto，运行时运维相关的 actuator.proto。例如 RPC 调用、发消息、读缓存、读配置这些都属于应用到运行时的，而 RPC  收请求、收消息、接收任务调度这些属于运行时调应用的，其它监控检查、组件管理、流量控制这些则属于运行时运维相关的。

这三个 proto 的示例可以看下图：

>![](https://gw.alipayobjects.com/mdn/rms_1c90e8/afts/img/A*J76nQoLLYWgAAAAAAAAAAAAAARQnAQ)

**运行时组件管控**

另外一方面，为了实现运行时的实现可替换，我们也在 MOSN 提了两个概念，我们把一个个分布式能力称为 Service，然后有不同的 Component 去实现这个 Service，一个 Service 可以有多个组件实现它，一个组件可以实现多个 Service。例如图里的示例就是有“MQ-pub” 这个发消息的 Service 有 SOFAMQ 和 Kafka 两个 Component 去实现，而 Kafka Component 则实现了发消息和健康检查两个 Service。
当业务真正通过 gRPC 生成的客户端发起请求的时候，数据就会通过 gRPC 协议发送给 Runtime，并且分发到后面一个具体的实现上去。这样的话，应用只需要使用同一套 API，通过请求里的参数或者运行时的配置，就对接到不同的实现。

>![](https://gw.alipayobjects.com/mdn/rms_1c90e8/afts/img/A*dK9rRLTvtlMAAAAAAAAAAAAAARQnAQ)

**运行时和 Mesh 的对比**

综上所述， 云原生应用运行时和刚才 Mesh 简单对比如下：

>![](https://gw.alipayobjects.com/mdn/rms_1c90e8/afts/img/A*xyu9T74SD9MAAAAAAAAAAAAAARQnAQ)

云原生应用运行时落地场景
从去年中开始研发，运行时目前在蚂蚁内部主要落地了下面几个场景。

**异构技术栈接入**

>![](https://gw.alipayobjects.com/mdn/rms_1c90e8/afts/img/A*8UJhRbBg3zsAAAAAAAAAAAAAARQnAQ)

在蚂蚁，不同的语言的应用除了 RPC 服务治理、消息等的需求之外，还希望使用上蚂蚁统一的中间件等基础设施能力，Java 和 Nodejs 是有对应的 SDK 的，而其他语言是没有的对应的 SDK 的。有了应用运行时之后，这些异构语言就可以直接通过 gRPC Client 调用运行时，对接上蚂蚁的基础设施。

**解除厂商绑定**

>![](https://gw.alipayobjects.com/mdn/rms_1c90e8/afts/img/A*eVoqRbkTFFwAAAAAAAAAAAAAARQnAQ)

刚才提到，蚂蚁的区块链、风控、智能客服、金融中台等等业务是既在主站部署，又有阿里云或者专有云部署的场景。有了运行时之后，应用可以一套代码和运行时一起出一个镜像，通过配置去决定调用哪个底层的实现，不跟具体的实现绑定。例如在蚂蚁内部对接的是 SOFARegistry 和 SOFAMQ 等产品，而到云上对接的是 Nacos、RocketMQ 等产品，到专有云对接的又是 Zookeeper、Kafka 等。这个场景我们正在落地当中。当然这个也可以用在遗留系统治理上，例如从 SOFAMQ 1.0 升级到 SOFAMQ 2.0，接了运行时的应用也无需升级。

**FaaS 冷启预热池**

FaaS 冷启预热池也是我们近期在探索的一个场景，大家知道 FaaS 里的 Function 在冷启的时候，是需要从创建 Pod 到下载 Function 再到启动的，这个过程会比较长。有了运行时之后，我们可以提前把 Pod 创建出来并启动好运行时，等到应用启动的时候其实已经非常简单的应用逻辑了，经过测试发现可以将从 5s 缩短 80% 到 1s。这个方向我们还会持续探索当中。

## 规划和展望

**API 共建**

运行时里最主要的一部分就是 API 的定义，为了落地内部，我们已经有一套较为完整的 API，但是我们也看到业界的很多产品有类似的诉求，例如 dapr、envoy 等等。所以接下来我们会去做的一件事情就是联合各个社区去推出一套大家都认可的云原生应用 API。

>![](https://gw.alipayobjects.com/mdn/rms_1c90e8/afts/img/A*d2BORogVotoAAAAAAAAAAAAAARQnAQ)

**持续开源**

另外我们近期也会将内部的运行时实践逐步开发出来，预计五六月份会发布 0.1 版本，并保持每月发布一个小版本的节奏，争取年底之前发布 1.0 版本。

>![](https://gw.alipayobjects.com/mdn/rms_1c90e8/afts/img/A*Kgr9QLc5TH4AAAAAAAAAAAAAARQnAQ)

## 总结

**最后做一下小结：**

1.Service Mesh 模式的引入是实现应用原云生的关键路径；

2.任何中间件兼可 Mesh 化，但研发效率问题依然部分存在；

3.Mesh 大规模落地是工程化的事情，需要完整的配套体系；

4.云原生应用运行时将是中间件等基础技术的未来形态，进一步解耦应用与分布式能力；

5.云原生应用运行时核心是 API，期望社区共建一个标准。

延伸阅读

- [带你走进云原生技术：云原生开放运维体系探索和实践](https://mp.weixin.qq.com/s?__biz=MzUzMzU5Mjc1Nw==&mid=2247488044&idx=1&sn=ef6300d4b451723aa5001cd3deb17fbc&chksm=faa0fdf6cdd774e03ccd9130099674720a81e7e109ecf810af147e08778c6582636769646490&scene=21)

- [积跬步至千里：QUIC 协议在蚂蚁集团落地之综述](https://mp.weixin.qq.com/s?__biz=MzUzMzU5Mjc1Nw==&mid=2247487717&idx=1&sn=ca9452cdc10989f61afbac2f012ed712&chksm=faa0ff3fcdd77629d8e5c8f6c42af3b4ea227ee3da3d5cdf297b970f51d18b8b1580aac786c3&scene=21)

- [Rust 大展拳脚的新兴领域：机密计算](https://mp.weixin.qq.com/s?__biz=MzUzMzU5Mjc1Nw==&mid=2247487576&idx=1&sn=0d0575395476db930dab4e0f75e863e5&chksm=faa0ff82cdd77694a6fc42e47d6f20c20310b26cedc13f104f979acd1f02eb5a37ea9cdc8ea5&scene=21)

- [Protocol Extension Base On Wasm——协议扩展篇](https://mp.weixin.qq.com/s?__biz=MzUzMzU5Mjc1Nw==&mid=2247487546&idx=1&sn=72c3f1ede27ca4ace7988e11ca20d5f9&chksm=faa0ffe0cdd776f6d17323466b500acee50a371663f18da34d8e4cbe32304d7681cf58ff9b45&scene=21)