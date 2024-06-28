# Ant Cloud Native Apps Exploring and Practice - ArchiSummit

> The introduction of the Mesh model is a key path to the application of clouds and ant groups have achieved mass landings internally.The sinking of more middleware capabilities, such as Message, DB, Cache Mesh and others, will be the future shape of intermediate technology when the app evolves from Mesh.Apps run to help developers construct cloud native apps quickly and to further decouple apps and infrastructure, while the app runs at the core of API standards, the community is expected to build together.

>![](https://gw.alipayobjects.com/mdn/rms_1c90e8/afts/img/A*nergRo8-RI0AAAAAAAAAAAAAARQnAQ)

## Ant Group Mesh Introduction

Ant is a technology and innovation-driven company, from its earliest days as a payment app on Taobao to its current services
As a large company with 1.2 billion users worldwide, Ant's technical architecture evolution will probably be divided into the following stages:

Prior to 2006, the earliest payment was a centralized monolithic application with modular development of different businesses.

In 2007, as more scenes of payments were promoted, an application and data splitting began to be made and some modifications to SOA were made.

After 2010, rapid payments, mobile payments, support for two-eleven and balance jewels have been introduced, and users have reached the level of hundreds of millions, and the number of ant applications has grown, and ants have developed many full sets of microservice middleware to support ant operations;

In 2014, like the advent of more business formalities like rush flow, online payments and more scenes, higher requirements for ant availability and stability, ants supported LDC moderation in microservice intermediation, off-site support for business support, and elasticity scaling-up in mixed clouds that support bi-11 ultra-mass traffic.

In 2020, ant business was not only digital finance, but also the emergence of new strategies such as digital life and internationalization, which prompted us to have a more efficient technical structure that would allow the operation to run faster and more steadily, so ant ants were able to internalize a more popular concept of cloud-origin in the industry.

>![](https://gw.alipayobjects.com/mdn/rms_1c90e8/afts/img/A*KCSVTZWSf8wAAAAAAAAAAAAAARQnAQ)

The technical structure of ant can also be seen to evolve along with the business innovations of the company from centralization to SOA to microservices, believing that the classmates with microservices are well known and that the practice of microservices to clouds has been explored by ants themselves in recent years.

## Why to introduce Service Mesh

Since ant has a complete set of microservice governance intermediaries, why do you need to introduce Service Mesh?

>![](https://gw.alipayobjects.com/mdn/rms_1c90e8/afts/img/A*Sq7oR6eO2QAAAAAAAAAAAAAAARQnAQ)

The service framework for ant self-research is SOFARPC as an example of a powerful SDK that includes a range of capabilities such as discovery of services, routing, melting out streams, etc.In a basic SOFA(Javaa) app, business code integrates SOFARP's SDK, both running in a process.After the large scale of sunk microservice, we faced some of the following problems with：

**Upgrade cost**：SDK requires business code introduction. Each upgrade requires a change code to be published.Because of the large scale of applications, some major technological changes or safety problems are being repaired.It takes thousands of apps to upgrade each time it takes time.
**Version Fragment**：is highly fragmented, due to the high cost of upgrades, which makes it difficult for us to use historical logic when writing our code and to evolve across technology.
**Cross-language is unmanageable**：ant online applications mostly use Java as a technical stack, but there are many cross-language applications in the front office, AI, Big Data, for example C++/Python/Golang etc. Their service governance capacity is missing due to SDK without a corresponding language.

We note that some concepts of Service Mesh in the cloud are beginning to emerge, so we are beginning to explore this direction.In the concept of Service Mesh, there are two concepts, one Control Plane Control and one Data Plane Dataplane.The core idea of the data plane is to decouple and to abstract some of the unconnected and complex logic (such as service discovery in RPC calls, service routing, melting breaks, security) into an independent process.As long as there is no change in the communications agreement between the operational and independent processes, the evolution of these capabilities can follow the autonomous upgrading of this independent process and the evolution of the entire Mesh can take place in a unified manner.Our cross-language applications, as long as the traffic passes through our Data Plane, are able to enjoy the capacities related to the governance of the services just mentioned, and the application of infrastructure capabilities to the bottom is transparent and truly cloud.

## Ant Mesh landing process

So, starting at the end of 2017, ant began to explore the technical direction of Service Mesh and presented a vision of a unified infrastructure with a sense of business upgrade.The main milestone is：

The Technology Advance Research Service Mesh technology was launched at the end of 2017 and set the direction for the future;

Beginning in early 2018 with Golang Self Research Sidecar MOSN and its source, mainly supporting RPC on a two-decimal scale pilot;

2019 New Message Mesh and DB Mesh shape in 618, covering a number of core links and exponentially 618

Two-11 years in 2019, covering hundreds of applications from all high-profile core links, supporting the Big Eleven at that time;

Twenty and eleven years in 2020, more than 80% of online applications are connected to the Mesh system and can be upgraded from capacity development to full capacity for 2 months.

## Ant Mesh Landing Architecture

Mesh at ant landing size is about thousands of applications and hundreds of thousands of levels of containers, a scale that falls in industry to a few and two times without a previous path to learn, so as ant arrives in a complete system of research and development delivery to support the mesh of ants as he arrives.

>![](https://gw.alipayobjects.com/mdn/rms_1c90e8/afts/img/A*eAlMT7SMTpMAAAAAAAAAAAAAARQnAQ)
Ant Mesh structure is probably our control plane, as shown in the graph, and the service end of the service governance centre, PaaS, monitoring centre, etc. are deployed as some of the existing products.There are also our transport systems, including R&D platforms and PaaS platforms.The middle is our main player data plane MOSN, which manages RPC, messages, MVC, Tasks four streams, as well as basic capabilities for health screening, monitoring, configuration, security, and technical risks, and MOSN blocks some interaction between operations and basic platforms.DBMesh is an independent product in the ant and is not drawn in the graph.Then the top tier is some of our applications that currently support access to many languages such as Java, Nodejs.
For applications, while infrastructure decoupling, access will require an additional upgrade cost, so in order to promote access to the app, ant makes the entire research and development delivery process, including by making the simplest access to the existing framework, by pushing forward in batches to manage risks and progress, and by allowing new applications default access to Mesh to do so.

At the same time, as sincerity grows, each of the capacities faced some problems of collaboration in R&D, and even of mutual impact on performance and stability, so that for the development effectiveness of the Mesh itself, we have made improvements in modular isolation, dynamic plugging of new capacities, automatic regression, and so on, which can be completed within two months from development to roll-out across the site.

## Explore on Cloud Native Apps Run

**New issues and reflections on mass backwardness**

Ant Mesh has now encountered some new problems with：
cross-language SDK maintenance master：Canada RPC examples. Most of the logic is already sinking into MOSN, but there is still some communication decoding protocol logic in Java, this SDK has some maintenance costs, how many lightweight SDKs, how many languages a team cannot have research and development in all languages. The quality of the Institute's code in this lightweight SDK is a problem.

A part of the application of the new：ant in business compatible with different environments is deployed both inside the ant and externally exported to financial institutions.When they are deployed to ant the control face of the ant and when the bank is received, the control of the bank is already in place.Most of the applications now contain a layer of their code and temporarily support the next when they meet unsupported components.

The earliest scenes from Service Mesh to Multi-Mesh：ant are Service Mesh, MOSN intercept traffic through network connecting agents, and other intermediates interact with the server through the original SDK.Now MOSN is more than a Service Mosh, but multi-Mesh, because, with the exception of RPC, we have supported more mesh Mesh landing sites, including messages, configurations, caches, etc.Each sinking intermediate can be seen, and almost all have a lightweight SDK on the side of the app, which, in the context of the first issue just a moment ago, finds a very large amount of lightweight SDK that needs to be maintained.In order to keep the features do not interact with each other, each feature opens different ports, calls with MOSN via different protocol.e.g. RPC protocol for RPC, MQ protocol for messages, cached Redis protocol.Then the current MOSN is more than just a flow orientation. For example, the configuration is to expose the API to use business code.

>![](https://gw.alipayobjects.com/mdn/rms_1c90e8/afts/img/A*80o8SYwyHJoAAAAAAAAAAAAAARQnAQ)
 
To solve the problems and scenes we are thinking about the following points：

Can the SDK be styled in different intermediaries, languages and languages?

Can interoperability protocols be unified?

3. Do we sink under our intermediate part to components or capabilities?

Can the implementation of the bottom be replaced?

>![](https://gw.alipayobjects.com/mdn/rms_1c90e8/afts/img/A*hsZBQJg0VnoAAAAAAAAAAAAAARQnAQ)

## Ant Cloud Native Apps Runtime Structure

Beginning last March, following several rounds of internal discussions and research into new ideas in industry, we introduced a concept of “cloud native apps” (hereinafter referred to as running on).By definition, we want this operation to include all distributive capabilities that the app cares for, help developers build your cloud native apps quickly, help apps and infrastructure to decouple more!

>![](https://gw.alipayobjects.com/mdn/rms_1c90e8/afts/img/A*iqQoTYAma4YAAAAAAAAAAAAAARQnAQ)

The core points of runtime design for cloud-native applications are as follows:

\*\*First \*\*, due to experience of MOSN sizing and associated shipping systems, we decided to build up our cloud native app on the basis of MOSN kernel.

\*\*Second \*\*, Abilities instead of Component Orientation, define the APIs for this running time.

**Third**, the interaction between business code and the Runtime API uses a uniform gRPC protocol so that the side of the business can generate a client directly and directly call through proto file.

**Four**'s component implementation after ability is replacable, for example, registration service provider may be SOFARegistry, or Nacos or Zookeper.

**Running abstract capabilities**

>![](https://gw.alipayobjects.com/mdn/rms_1c90e8/afts/img/A*hWIVR6ccduYAAAAAAAAAAAAAARQnAQ)

To abstract some of the capabilities most needed for cloud apping, we set a few principles：

1. Follow the APIs and Scenarios required for distributed apps instead of components;
   2.APIs are intuitive, used in boxes, and are better than configured;
   3.APIs are not bound to implement and differentiate using extension fields.

With this principle, we abstract out the primary API, which is the app for mosn.proto, the appcallback.proto for the app when running, and the relevant actuator.proto for the app when running.For example, RPC calls, messages, read caches, read configurations are all applied to running, while RPC receipts, messages, incoming task schedules, are applied when running. Other control checks, component management, traffic controls are related to running wikes.

Three examples of this proto can be seen at：

>![](https://gw.alipayobjects.com/mdn/rms_1c90e8/afts/img/A*J76nQoLLYWgAAAAAAAAAAAAAARQnAQ)

**Run Component Controls**

On the other hand, we have two concepts in MOSN for the purpose of realizing replaceability when running. We call a distribution capability and then have a different component to perform this Service, a service that can be implemented with multiple components, and a component that can deliver multiple services.For example, the example in the graph is that the service with the message "MQ-pub" is implemented by SOFAMQ and Kafka Component, while Kafka Component implements both the message and health check service.
When a transaction is actually requested via a gRPC-generated client, the data will be sent to Runtime via the gRPC protocol and distributed to the next specific implementation.In this way, the app needs to use only the same set of API, which can be implemented differently by the parameters in the request or when the configuration is running.

>![](https://gw.alipayobjects.com/mdn/rms_1c90e8/afts/img/A*dK9rRLTvtlMAAAAAAAAAAAAAARQnAQ)

**Compare between runtime and Mesh**

Based on the above, when the cloud app is running and just just Mesh are easy to compare with：

>![](https://gw.alipayobjects.com/mdn/rms_1c90e8/afts/img/A*xyu9T74SD9MAAAAAAAAAAAAAARQnAQ)

Scene
started research last year while the cloud native app is running. The following scenes are currently falling inside the ant area.

**Isomer Technical Stack Access**

>![](https://gw.alipayobjects.com/mdn/rms_1c90e8/afts/img/A*8UJhRbBg3zsAAAAAAAAAAAAAARQnAQ)

In the case of ants, applications in different languages, in addition to the need for RPC service governance, messages, etc., the infrastructure capabilities such as the one-size-fits-all intermediate of the ant are desirable and Java and Nodejs have corresponding SDKs, while the other languages are not corresponding SDKs.After the application runs, these isomer languages can be used directly through GRPC Client to the ant infrastructure.

**Unbind the manufacturer**

>![](https://gw.alipayobjects.com/mdn/rms_1c90e8/afts/img/A*eVoqRbkTFFwAAAAAAAAAAAAAARQnAQ)

As mentioned earlier, ant blockchains, wind control, intelligent support services, financial intermediaries, etc., are scenes where they are deployed on their main stations, where there is either Aliyun or cloud.After running, the app can combine a set of code with a mirror when running. By configuring it to determine which bottom layer of implementation to be called, without being bound to specific implementations.For example, the internal interface between ant is for products such as SOFARegistration and SOFAMQ, and on the cloud is for products such as Nacos, RocketMQ, to Zokeper, Kafka and others.This scenario is in the process of reaching us.Of course, this can also be used for legacy system governance, such as upgrading from SOFAMQ 1.0 to SOFAMQ 2.0, and then running apps need not be upgraded.

\*\*FaaS Cold Pool Pool \*\*

FaaS Cool is also a recent scene we are exploring and you know that the Function in FaaS needs to go from Pod creation to Download Function to Start, a process that will be lengthy.After running time, we can create Pod in advance and start up good running. Wait a very simple app logic when the app starts. Test it can be shortened from 5s to 1s.We will continue to explore this direction as well.

## Planning and outlook

**API**

The most important part of the running time is the definition of the API. We already have a more complete set of APIs for the sake of getting inside, but we also see that many products in industry have similar demands, such as dapr, envoy, etc.So one of the next things we will do is to bring together communities to launch a set of recognized cloud native API.

>![](https://gw.alipayobjects.com/mdn/rms_1c90e8/afts/img/A*d2BORogVotoAAAAAAAAAAAAAARQnAQ)

**Continuous Open Source**

We will also develop our internal running practice in the near future, with a release of 0.1 in May and June, and keep a small monthly release pace, aiming to publish 1.0 by the end of the year.

>![](https://gw.alipayobjects.com/mdn/rms_1c90e8/afts/img/A*Kgr9QLc5TH4AAAAAAAAAAAAAARQnAQ)

## Summary

**Last Summary：**

1.Service Mesh mode introduction is a key path to the application of the cloud;

Any mesh that allows Mesh to be generated, but the problem of R&D efficiency remains partially present;

3.Mesh Large-scale landfall is a matter of engineering and requires a complete suite of systems;

4. Cloud native applications will be the future shape of basic technologies such as intermediaries, further decoupling and distributive capabilities;

The cloud native app runs at the heart of the API, and the community is expected to build one standard together.

Extend Reading

- [Take you into Cloud Native Technology：Native Open Delivery Systems Exploration and Practices](https://mp.weixin.qq.com/s?_biz=MzUzU5Mjc1Nw===\&mid=2247488044\&idx=1\&sn=e6300d4b451723a5001cd3deb17fbc\&chksm=faa0f6cdd774e03ccd91300996747a8e7e109ecf810af147e08c663676946490\&scene=21)

- [Taking a thousand miles one step at a time: A comprehensive overview of the QUIC protocol landing at Ant Group](https://mp.weixin.qq.com/s?__biz=MzUzMzU5Mjc1Nw==\&mid=2247487717\&idx=1\&sn=ca9452cdc10989f61afbac2f012ed712\&chksm=faa0ff3fcdd77629d8e5c8f6c42af3b4ea227ee3da3d5cdf297b970f51d18b8b1580aac786c3\&scene=21)

- [Rust's emerging field showing its prowess: confidential computing](https://mp.weixin.qq.com/s?__biz=MzUzMzU5Mjc1Nw==\&mid=2247487576\&idx=1\&sn=0d0575395476db930dab4e0f75e863e5\&chksm=faa0ff82cdd77694a6fc42e47d6f20c20310b26cedc13f104f979acd1f02eb5a37ea9cdc8ea5\&scene=21)

- [Protocol Extension Base on Wasm — protocol extension](https://mp.weixin.qq.com/s?_biz=MzUzU5Mjc1Nw===\&mid=2247487546\&idx=1\&sn=72c3f1e27ca4ace788e11ca20d5f9\&chksm=faa0ffe0cd776f6d17323466b500acee50a371663f18da34d8e4d72304d7681cf589b45\&scene=21)
