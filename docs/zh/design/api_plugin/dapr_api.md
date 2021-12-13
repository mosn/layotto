## TLDR
本提案想让开源Layotto 既支持 Layotto API，又支持Dapr API。类似于“Minio 既支持Minio API，又支持AWS S3 API”

## 想解决的问题
1. 目前，我们尽量保证Layotto API 里各个字段定义的和Dapr一样，但用户真正关心的是两者sdk能否复用。虽然我们努力保证proto字段一致，但只要不能复用sdk就没解决用户的问题，还给自己增加维护成本。
   比如：
   ![image](https://user-images.githubusercontent.com/26001097/145837477-00fc5cd8-32eb-4ce9-bbfb-6e590172fce8.png)

因此，我们想让Layotto直接支持Dapr的grpc API （一模一样，包括 package名），对于用户来说，他可以用Dapr sdk在两者之间自由切换，不用担心被厂商绑定。

2. 另一方面，还需要有一定扩展性。我们在生产落地的过程中发现目前的Dapr API没法完全满足需求，难免要对API做一些扩展。扩展的API已经加到了现在的 Layotto API里，提案提给了 Dapr 社区、但还在等社区慢慢接受，比如config API，比如Lock API

## 方案
### Layotto API on Dapr API
![image](https://user-images.githubusercontent.com/26001097/145838604-e3a0caad-9473-4092-a2c6-0cc46c972790.png)
1. Layotto 会启动一个grpc服务器，前阵子刚加了个API插件功能，我们可以通过 API插件的形式注册一个Dapr API 插件；
2. 另一方面，保留Layotto API。Layotto接收到Layotto API请求后，翻译成Dapr API，然后按Dapr API处理。
   这样的好处有：
- 复用代码
- Layotto API 可以按生产需求进行扩展，比如支持Lock API，configuration API等；做了扩展后可以提给Dapr 社区，再慢慢讨论，即使最终讨论的结果和原始方案不一样，也只是影响最终做出来的Dapr API，不会影响已经用上Layotto API的用户。

### 用户价值
对于用户来说：

- 如果用户担心厂商绑定，可以只用Dapr API，可以用同一套Dapr sdk在Dapr 和Layotto 之间迁移，减少用户疑虑；
- 如果用户相信我们的落地经验、愿意用Layotto API，那他们可以用Layotto API，代价是没法用同一个sdk在两个sidecar之间迁移

### Q&A
#### 想给Dapr API加字段怎么加
##### 想加个字段(field)
比如想给layotto api加个abc字段，可以通过metadata或者grpc头把这个字段传给dapr API
dapr API的实现再把这个字段透传给组件，组件解析这个字段

##### 不只加字段，还要加一些逻辑、机制（mechanism）
比如layotto api加个abc字段，如果abc==true，那么runtime走一段特殊逻辑

这种情况要修改Dapr API的实现，加一段if else

#### 想加新API怎么加
加在layotto API上，新API不需要复用Dapr API；等Dapr接收提案后再修改实现，layotto API不变
