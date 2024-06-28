# MOSN subproject Layotto：opens the service grid + new chapter when app runs

> Author profile：
> Magnetic Army. Fancy is an ancient one, cultivating for many years in the infrastructure domain, with in-depth practical experience of Service Mosh, and currently responsible for the development of projects such as MOSN, Layotto and others in the middle group of ant groups.
> Layotto official GitHub address: [https://github.com/mosn/layotto](https://github.com/mosn/layotto)

Click on a link to view the live video：[https://www.bilibili.com/video/BV1hq4y1L7FY/](https://www.bilibili.com/video/BV1hq4y1L7FY/)

Service Mesh is already very popular in the area of microservices, and a growing number of companies are starting to fall inside, and ants have been investing heavily in this direction from the very beginning of Service Mesh programme. So far, the internal Mesh programme has covered thousands of applications, hundreds of thousands of containers and has been tested many times, the decoupling of business coupling brought about by Service Mosh, smooth upgrades and other advantages have greatly increased iterative efficiency in intermediaries.

We have encountered new problems after mass landings, and this paper focuses on a review of service Mesh's internal landings and on sharing solutions to new problems encountered after service Mesh landing.

## Service Mesh Review and Summary

### Instrument for standardized international reporting of military expenditures

> ![](https://gw.alipayobjects.com/mdn/rms_1c90e8/afts/img/A*p8tGTbpLRegAAAAAAAAAAAAAARQnAQ)
> Under the microservice architecture, infrastructure team typically provides a SDK that encapsulates the ability to govern the various services, while ensuring the proper functioning of the application, it is also clear that each infrastructure team iterates a new feature that requires the involvement of the business party to use it, especially in the bug version of the framework, often requiring a forceful upgrade of the business side, where every member of the infrastructure team has a deep sense of pain.

The difficulties associated with upgrading are compounded by the very different versions of the SDK versions used by the application and the fact that the production environment runs in various versions of the SDK, which in turn makes it necessary to consider compatibility for the iterations of new functions as if they go ahead with the shacks, so that the maintenance of the code is very difficult and some ancestral logic becomes uncareful.

The development pattern of the “heavy” SDKs makes the governance of the isomer language very weak and the cost of providing a functionally complete and continuously iterative SDK for all programming languages is imaginable.

In 18 years, Service Mesh continued to explode in the country, a framework concept designed to decouple service governance capacity with business and allow them to interact through process-level communications.Under this architecture model, service governance capacity is isolated from the application and operated in independent processes, iterative upgrading is unrelated to business processes, which allows for rapid iterations of service governance capacity, and each version can be fully upgraded because of the low cost of upgrading, which has addressed the historical burden and the SDK “light” directly reduces the governance threshold for isomer languages and no longer suffers from the SDK that needs to develop the same service governance capability for each language.

### Current status of service Mesh landings

> ![](https://gw.alipayobjects.com/mdn/rms_1c90e8/afts/img/A*rRG_TYlHMqYAAAAAAAAAAAAAARQnAQ)
> 蚂蚁很快意识到了 Service Mesh 的价值，全力投入到这个方向，用 Go 语言开发了 MOSN 这样可以对标 envoy 的优秀数据面，全权负责服务路由，负载均衡，熔断限流等能力的建设，大大加快了公司内部落地 Service Mesh 的进度。

Now that MOSN has overwritten thousands of apps and hundreds of thousands of containers inside ant ants, newly created apps have default access to MOSN to form closers.And MOSN handed over a satisfactory： in terms of resource occupancy and loss of performance that is of greatest concern to all.

1. RT is less than 0.2 ms

2. Increase CPU usage by 0% to 2%

3. Memory consumption growth less than 15M

The technical stack of the NodeJS, C+++ isomers is also continuously connected to MOSN due to the Service Mesh service management thresholds that lower the isomer language.

After seeing the huge gains from RPC capacity Mih, internal ants also transformed MQ, Cache, Config and other middleware capabilities, sinking to MOSN, improving the iterative efficiency of the intermediate product as a whole.

### C. New challenges

1. Apply strong binding to infrastructure

> ![](https://gw.alipayobjects.com/mdn/rms_1c90e8/afts/img/A*nKxcTKLp4EoAAAAAAAAAAAAAARQnAQ)
> A modern distributed application often relies on RPC, Cache, MQ, Config and other distributed capabilities to complete the processing of business logic.

When RPC was initially seen, other capabilities were quickly sinking.Initially, they were developed in the most familiar way, leading to a lack of integrated planning management, as shown in the graph above, which relied on SDKs of a variety of infrastructure, and in which SDK interacted with MOSN in a unique way, often using private agreements provided by the original infrastructure, which led directly to a complex intermediate capability, but in essence the application was tied to the infrastructure, such as the need to upgrade the SDK from Redis to Memcache, which was more pronounced in the larger trend of the application cloud, assuming that if an application was to be deployed on the cloud, because the application relied on a variety of infrastructures, it would be necessary to move the entire infrastructure to the cloud before the application could be successfully deployed.
So how to untie the application to the infrastructure so that it can be transplantable and that it can feel free to deploy across the platform is our first problem.

2. Isomal language connectivity

> ![](https://gw.alipayobjects.com/mdn/rms_1c90e8/afts/img/A*oIdQQZmgtyUAAAAAAAAAAAAAARQnAQ)
> It has been proved that Service Mesh does reduce the access threshold of heterogeneous languages, but after more and more basic capabilities sink to MOSN, we gradually realized that in order to allow applications to interact with MOSN, various SDKS need to develop communication protocols and serialization protocols. If you add in the need to provide the same functionality for a variety of heterogeneous languages, the difficulty of maintenance increases exponentially

Service Mesh has made the SDK historic, but for the current scenario of programming languages and applications with strong infrastructural dependence, we find that the existing SDK is not thin enough, that the threshold for access to the isomer language is not low enough and that the threshold for further lowering the isomer language is the second problem we face.

## Multi Runtime Theory Overview

### A, what is Runtime?

> ![](https://gw.alipayobjects.com/mdn/rms_1c90e8/afts/img/A*hQT-Spc5rI4AAAAAAAAAAAAAARQnAQ)
> In the early 20th century, Bilgin lbryam published a paper called
> Multi-Runtime Microservices Architecture
> This article discusses the shape of the next phase of microservices architecture.

As shown in the graph above, the author abstracts the demand for distributed services and is divided into four chaos：

1. Life Cycle (Lifecycle)
   mainly refers to compilation, packing, deployment and so forth, and is largely contracted by docker and kubernetes in the broad cloud of origins.

2. Network (Networking)
   A reliable network is the basic guarantee of communication between microservices, and Service Mesh is trying to do so and the stability and usefulness of the current popular data face of MOSN and envoy have been fully tested.

3. The status (State)
   services that are required for distribution systems, workflow, distribution single, dispatching, power equivalent, state error restoration, caching, etc. can be uniformly classified as bottom status management.

4. Binding (Binding)
   requires not only communication with other systems but also integration of various external systems in distributed systems, and therefore has strong reliance on protocol conversion, multiple interactive models, error recovery processes, etc.

After the need has been clarified, drawing on the ideas of Service Mesh, the author has summarized the evolution of the distributed services architecture as： below.

> ![](https://gw.alipayobjects.com/mdn/rms_1c90e8/afts/img/A*rwS2Q5yMp_sAAAAAAAAAAAAAARQnAQ)
> Phase I is to decouple infrastructure capabilities from the application and to convert them into an independent residecar model that runs with the application.

The second stage is to unify the capabilities offered by the sidecar into a single settlement run from the development of the basic component to the development of the various distributive capabilities to the development of the various distributive capacities, completely block the details of the substrate and, as a result of the ability orientation of the API, the application no longer needs to rely on SDK from a wide range of infrastructures, except for the deployment of the APIs that provide the capabilities.

The author's thinking is consistent with what we want to resolve, and we have decided to use the Runtime concept to solve the new problems that Service Mesh has encountered to date.

### B, Service Mesh vs Runtime

> ![](https://gw.alipayobjects.com/mdn/rms_1c90e8/afts/img/A*srPVSYTEHc4AAAAAAAAAAAAAARQnAQ)
> In order to create a clearer understanding of Runtime, a summary of Service Mesh with regard to the positioning, interaction, communication protocols and capacity richness of the two concepts of Runtime is shown, as can be seen from Service Mosh, when Runtime provides a clearly defined and capable API, making the application more straightforward to interact with it.

## MOSN sub-project Layotto

### A, dapr research

> ![](https://gw.alipayobjects.com/mdn/rms_1c90e8/afts/img/A*Ab9HSYIK7CQAAAAAAAAAAAAAARQnAQ)
> dapr is a well-known Runtime product in the community and has a high level of activity, so we first looked at the dapr case, finding that the dapr has the following advantage of：

1. A variety of distributive capabilities are provided, and the API is clearly defined and generally meets the general usage scenario.

2. Different delivery components are provided for each capability, essentially covering commonly used intermediate products that can be freely chosen by users as needed.

When considering how to set up a dapr within a company, we propose two options, such as the chart： above

1. Replace：with the current MOSN and replace with the dapr. There are two problems with：

Dapr does not currently have the full range of service governance capabilities included in Service Mesh although it provides many distributive capabilities.

b. MOSN has fallen on a large scale within the company and has been tested on numerous occasions with the direct replacement of MOSN stability by a dapr.

2. In：, add a dapr container that will be deployed with MOSN in two sidecar mode.This option also has two problems with：

The introduction of a new sidecar will require consideration of upgrading, monitoring, infusion and so forth, and the cost of transport will soar.

b. The increased maintenance of a container implies an additional risk of being hacked and this reduces the availability of the current system.

Similarly, if you are currently using envoy as a data face, you will also face the above problems.
We therefore wish to combine Runtime with Service Mesh and deploy through a full sidecar to maximize the use of existing MSh capabilities while ensuring stability and the constant cost of delivery.In addition, we hope that, in addition to being associated with MOSN, the capacity of the RPF will be combined in the future with envoy to solve the problems in more scenarios, in which Layotto was born.

### Layout B & Layout

> ![](https://gw.alipayobjects.com/mdn/rms_1c90e8/afts/img/A*sdGoSYB_XFUAAAAAAAAAAAAAARQnAQ)
> As shown in the above chart, Layotto is above all over the infrastructure and provides a standard API for upper-tier applications with a uniform range of distributive capabilities.For Layotto applications, developers no longer need to care for differences in the implementation of substrate components, but just what competencies the app needs and then call on the adaptive API, which can be completely untied to the underlying infrastructure.

For applications, interaction is divided into two blocks, one as a standard API for GRPC Clients calling Layotto and another as a GRPC Server to implement the Layotto callback and benefit from the gRPC excellent cross-language support capability, which no longer requires attention to communications, serialization, etc., and further reduces the threshold for the use of the technical stack of isomers.

In addition to its application-oriented, Layotto also provides a unified interface to the platform that feeds the app along with the sidecar state of operation, facilitates SRE peer learning to understand the state of the app and make different initiatives for different states, taking into account existing platform integration with k8s and so we provide access to HTTP protocol.

In addition to Layotto itself design, the project involves two standardized constructions, firstly to develop a set of terminological clocks; the application of a broad range of APIs is not an easy task. We have worked with the Ari and Dapr communities in the hope that the building of the Runtime API will be advanced, and secondly for the components of the capabilities already achieved in the dapr community, our principle is to reuse, redevelop and minimize wasting efforts over existing components and repeat rotations.

In the end, Layotto is now built over MOSN, we would like Layotto to be able to run on envoy, so that you can increase Runtime capacity as long as you use Service Mesh, regardless of whether the data face is used by MOSN or envoy.

### C, Layotto transplantation

> ![](https://gw.alipayobjects.com/mdn/rms_1c90e8/afts/img/A*2DrSQJ6GL8cAAAAAAAAAAAAAARQnAQ)
> As shown in the graph above, once the standardisation of the Runtime API is completed, access to Layotto applications is naturally portable, applications can be deployed on private clouds and various public clouds without any modification, and since standard API is used, applications can be freely switched between Layotto and dapr without any modification.

### Meaning of name

> ![](https://gw.alipayobjects.com/mdn/rms_1c90e8/afts/img/A*CCckTZ_gRsMAAAAAAAAAAAAAARQnAQ)
> As can be seen from the above schematic chart, the Layotto project itself is intended to block the details of the infrastructure and to provide a variety of distributive capabilities to the upper level of application. This approach is as if it adds a layer of abstraction between the application and the infrastructure, so we draw on the OSI approach to defining a seven-tiered model of the network and want Layot to serve the eighth tier of the application, to be 8 in Italian, Layer otto is meant to simplify to become Layotto, along with Project Code L8, which is also the eighth tier and is the source of inspiration for

An overview of the completion of the project is presented below, with details of the achievement of four of its main functions.

### E. Configuration of original language

> ![](https://gw.alipayobjects.com/mdn/rms_1c90e8/afts/img/A*mfkRQZH3oNwAAAAAAAAAAAAAARQnAQ)
> First is the configuration function commonly used in distributed systems, applications generally use the configuration center to switch or dynamically adjust the running state of the application.The implementation of the configuration module in Layotto consists of two parts. One is a reflection on how to define the API for this capability, and one is a specific implementation, each of which is seen below.

It is not easy to define a configuration API that meets most of the actual production demands. Dapr currently lacks this capability, so we worked with Ali and the Dapr community to engage in intense discussions on how to define a version of a reasonable configuration API.

As the outcome of the discussions has not yet been finalized, Layotto is therefore based on the first version of the draft we have submitted to the community, and a brief description of our draft is provided below.

We first defined the basic element： for general configuration

1. appId：indicates which app the configuration belongs to

2. Key configured for key：

3. Value of content：configuration

4. group：configurations are configured. If an appId is too many configurations, we can group these configurations for maintenance.

In addition, we added two advanced features to suit more complex configurations using Scene：

1. label, used to label configurations, such as where the configuration belongs, and when conducting configuration queries, we'll use label + key to query configuration.

2. tags, users give configuration additional information such as description, creator information, final modification time, etc. to facilitate configuration management, audits, etc.

For the specific implementation of the configuration API as defined above, we currently support query, subscription, delete, create, and modify five kinds of actions in which subscriptions to configuration changes use the stream feature of GRPC and the components where the configuration capacity is implemented at the bottom, we have selected the domestically popular apollo and will add others later depending on demand.

### F: Pub/Sub

> ![](https://gw.alipayobjects.com/mdn/rms_1c90e8/afts/img/A*YJs-R6WFhkgAAAAAAAAAAAAAARQnAQ)
> for Pub/Sub capabilities, we have explored the current implementation of dapr and have found that we have largely met our needs, so we have directly reintroduced the Dapr API and components that have been suitably matched in Layotto, which has saved us a great deal of duplication and we would like to maintain a collaborative approach with the dapr community rather than repeat the rotation.

Pub is an event interface provided by the App calls Layotto and the Sub function is one that implements ListTopicSubscriptions with OnTopicEvent in the form of a gRPC Server, one that tells Layotto apps that need to subscribe to which topics, and a callback event for Layotto receive a change in top.

Dapr for the definition of Pub/Sub basically meets our needs, but there are still shortfalls in some scenarios, dapr uses CloudEvent standards, so the pub interface does not return value, which does not meet the need in our production scenes to require pub messages to return to the messageID that we have already submitted the needs to dapr communities, and we are waiting for feedback, taking into account mechanisms for community asynchronous collaboration, we may first increase the results and then explore with the community a better compatibility programme.

### G and RPC original

> ![](https://gw.alipayobjects.com/mdn/rms_1c90e8/afts/img/A*i-JnSaeZbJ4AAAAAAAAAAAAAARQnAQ)
> The capacity of RPC is not unfamiliar and may be the most basic needs under the microservice architecture, the definition of RPC interfaces, we also refer to the dapr community definition and therefore the interface definition is fully responsive to our needs and thus the interface definition is a direct reuse of dapr but the current RPC delivery programme provided by dapr is still weak, and MOSN is very mature over the years, This is a brave combination of Runtime with Service Mesh and MOSN itself as a component of our capacity to implement RPC and thereby Layotto submit to MOSN for actual data transfer upon receipt of RPC requests, The option could change routing rules through istio, downgraded flow and so on, which would amount to a direct replication of Service Mesh's capabilities. This would also indicate that Runtime is not about listing the Service Mesh, but rather a step forward on that basis.

In terms of details, in order to better integrate with MOSN, we have added one Channel, default support for dubbo, bolt, HTTP three common RPC protocols to the RPC. If we still fail to meet the user scene, we have added Before/After filter to allow users to customize extensions and implement protocol conversions, etc.

### H, Actuator

> ![](https://gw.alipayobjects.com/mdn/rms_1c90e8/afts/img/A*E_Q-T4d_bm4AAAAAAAAAAAAAARQnAQ)
> In actual production environments, in addition to the various distributive capabilities required for the application, we often need to understand the operational state of the application, based on this need, we abstract an actuator interface, and we currently do not have the capability to do so at the moment, dapr and so we are designed on the basis of internal demand scension scenarios to expose the full range of information on the application at the startup and running stages, etc.

Layotto divides the exposure information into an individual：

1. Health：This module determines whether the app is healthy, e.g. a strongly dependent component needs to be unhealthy if initialization fails, and we refer to k8s for the type of health check to：

a. Readiness：indicates that the app is ready to start and can start processing requests.

b. Liveness：indicates the state of life of the app, which needs to be cut if it does not exist.

2. Info：This module is expected to expose some of the dependencies of the app, such as the service on which the app depends, the subscription configuration, etc. for troubleshooting issues.

Health exposure health status is divided into the following Atlash：

1. INIT：indicates that the app is still running. If the app returns this value during the release process, the PaaS platform should continue waiting for the app to be successfully started.

2. UP：indicates that the app is starting up normally, and if the app returns this value, the PasS platform can start loading traffic.

3. DOWN：indicates that the app failed to boot, meaning PaaS needs to stop publishing and notify the app owner if the app returns this value during the release process.

The search for Layotto is now largely complete in the Runtime direction, and we have addressed the current problems of infrastructure binding and the high cost of isomer language access using a standard interactive protocol such as gRPC to define a clearly defined API.As the future API standardises the application of Layotto can be deployed on various privately owned and publicly owned clouds, on the one hand, and free switching between Layotto, dapr and more efficient research and development, on the other.

Currently, Serverless fields are also flown and there is no single solution, so Layotto makes some attempts in Serverless directions, in addition to the input in the Runtime direction described above.

## Exploring WebAssembly

### Introduction to the Web Assembly

> ![](https://gw.alipayobjects.com/mdn/rms_1c90e8/afts/img/A*-ACRSpqbuJ0AAAAAAAAAAAAAARQnAQ)
> WebAssembly, abbreviated WASM, a collection of binary commands initially running on the browser to solve the JavaScript performance problems, but due to its good safety, isolation, and linguistic indifference, one quickly starts to get it to run outside the browser. With the advent of the WASI definition, only one WASM will be able to execute the WAS document anywhere.

Since WebAssembly can run outside the browser, can we use it in Serverless fields?Some attempts had been made in that regard, but if such a solution were to be found to be a real one, it would be the first question of how to address the dependence of a functioning Web Assembly on infrastructure.

### Principles of B and Web Assembly landing

Currently MOSN runs on MOSN by integrating WASM Runtime to meet the need for custom extensions to MOSN.Layotto is also built over MOSN so we consider combining the two in order to implement the following graph：

> ![](https://gw.alipayobjects.com/mdn/rms_1c90e8/afts/img/A*U7UDRYyBOvIAAAAAAAAAAAAAARQnAQ)
> Developers can develop their code using a variety of preferred languages such as Go/C+/Rust and then run them over MOSN to produce WASM files and call Layotto provide standard API via local function when WASM style applications need to rely on various distribution capabilities in processing requests, thereby directly resolving the dependency of WASM patterns.

Layotto now provides Go with the implementation of the Rust version of WASM, while supporting the demo tier function only, is enough for us to see the potential value of such a programme.

In addition, the WASM community is still in its early stages and there are many places to be refined, and we have submitted some PRs to the community to build up the backbone of the WASM technology.

### C. WebAssembly Landscape Outlook

> ![](https://gw.alipayobjects.com/mdn/rms_1c90e8/afts/img/A*NzwKRY2GZPcAAAAAAAAAAAAAARQnAQ)
> Although the use of WAS in Layotto is still in the experimental stage, we hope that it will eventually become a service unless it is developed through a variety of programming languages, as shown in the graph above, and then codify the WASM document, which will eventually run on Layotto+MOSN, while the application wiki management is governed by k8, docker, prometheus and others.

## Community planning

Finally, look at what Layotto does in the community.

### A, Layotto vs Dapr

> ![](https://gw.alipayobjects.com/mdn/rms_1c90e8/afts/img/A*OpQTRqoMpK0AAAAAAAAAAAAAARQnAQ)
> Charted Layotto in contrast to the existing capabilities in Layotto, our development process at Layotto, always aim to achieve the goal of a common building, based on the principle of re-use, secondary development, and for the capacity being built or to be built in the future, we plan to give priority to Layotto and then to the community to merge into standard API, so that in the short term it is possible that the Layotto API will take precedence over the community, but will certainly be unified in the long term, given the mechanism for community asynchronous collaboration.

### The APP

> ![](https://gw.alipayobjects.com/mdn/rms_1c90e8/afts/img/A*nKxcTKLp4EoAAAAAAAAAAAAAARQnAQ)
> We have had extensive discussions in the community about how to define a standard API and how Layotto can run on envoy, and we will continue to do so.

### C, Road map

> ![](https://gw.alipayobjects.com/mdn/rms_1c90e8/afts/img/A*iUV3Q7S3VLEAAAAAAAAAAAAAARQnAQ)
> Layotto currently support four major functionalities in support of RPC, Config, Pub/Sub, Actuator and is expected to devote attention to distribution locks and observations in September, and Layotto plugging in December, which it will be able to run on envoy, with the hope that further outputs will be produced for the WebCongress exploration.

### Official open source

> ![](https://gw.alipayobjects.com/mdn/rms_1c90e8/afts/img/A*S6mdTqAapLQAAAAAAAAAAAAAARQnAQ)
> gave a detailed presentation of the Layotto project and most importantly the project is being officially opened today as a sub-project of MOSN and we have provided detailed documentation and demo examples to facilitate quick experience.

The construction of the API standardization is a matter that needs to be promoted over the long term, while standardization means not meeting one or two scenarios, but the best possible fitness for most use scenarios, so we hope that more people can participate in the Layotto project, describe your use scenario, discuss the API definition options, come together to the community, ultimately reach the ultimate goal of Write once, Run any!
