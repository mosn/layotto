# Multi-Runtime 2022：待解决的问题
## 1. API标准建设
根据落地用户的生产需求继续建设API标准、提交给Dapr社区共建。比如：
- 分布式锁API
- 配置API
- 延迟消息API

## 2. 生态共建
如何让已经落地Service Mesh的用户平滑迁移到Multi-Runtime？目前在做的一件事是Layotto on Envoy支持;

能否让Runtime API更好的融入K8S生态？目前在做的事是Layotto集成进k8s生态;


## 3. 服务早期生产用户
开源要做通用的、解决生产问题的功能。观察早期生产用户，目前面临以下问题：
### 3.1. 扩展性
让整个项目可扩展，比如某个公司想用layotto但是又想扩展一些自己的功能，要么能自己起一个项目、import开源layotto后通过钩子做一些扩展，要么能通过动态连接库之类的办法去扩展layotto二进制文件。目前这两种办法,dapr和layotto都没法做到，想扩展只能fork出来改代码
### 3.2. 稳定性风险
import开源Layotto之后，panic风险巨大,因为依赖了Dapr所有组件，这些组件用的库五花八门，可能panic，可能依赖冲突。能否通过按需编译、隔离性设计来减少panic风险？

目前开源项目的测试投入相对于公司里的测试流程来说少太多了，怎么建设开源测试体系；

### 3.3. 可观测性
> 以前没service mesh的时候，有问题我能自己查；后来有了service mesh，遇到问题我只能找别人来查了
> ——某测试同学

在生产环境落地Service Mesh会导致排查问题变难，而 Multi-Runtime 下沉的功能多了，排查起来更难。
要建设 Multi-Runtime 可观测性，避免让生产用户查问题难上加难。

## 4. 新研发模式
sidecar 支持 serverless 落地；