# 源码解析 7层流量治理,接口限流

>作者简介：
>张晨，是开源社区的爱好者，致力于拥抱开源，希望能和社区的各位开源爱好者互相交流互相进步和成长。
> 
>写作时间: 2022年4月20日


## Overview
此文档的目的在于分析接口限流的实现

## 前提：
文档内容所涉及代码版本如下

[https://github.com/mosn/mosn](https://github.com/mosn/mosn)

Mosn   d11b5a638a137045c2fbb03d9d8ca36ecc0def11（develop分支）

## 源码分析
### 总体分析
参考 <br />[https://mosn.io/docs/concept/extensions/](https://mosn.io/docs/concept/extensions/)

Mosn 的 Stream Filter 扩展机制

![01.png](https://gw.alipayobjects.com/mdn/rms_5891a1/afts/img/A*tSn4SpIkAa4AAAAAAAAAAAAAARQnAQ)

### 代码均在： [flowcontrol代码](https://github.com/mosn/mosn/tree/master/pkg/filter/stream/flowcontrol)

### stream_filter_factory.go分析
此类为一个工厂类，用于创建 StreamFilter

定义了一些常量用作默认值

![02.png](https://gw.alipayobjects.com/mdn/rms_5891a1/afts/img/A*PAWCTL6MS40AAAAAAAAAAAAAARQnAQ)

定义了限流配置类用作加载yaml定义并且解析生产出对应的功能

![03.png](https://gw.alipayobjects.com/mdn/rms_5891a1/afts/img/A*Ua32SokhILEAAAAAAAAAAAAAARQnAQ)

init() 初始化内部就是将 name 和 对应构造函数存储到 filter拦截工厂的map中

![04.png](https://gw.alipayobjects.com/mdn/rms_5891a1/afts/img/A*kb3qRqWnqxYAAAAAAAAAAAAAARQnAQ)

着重讲一下 createRpcFlowControlFilterFactory  生产出rpc流控工厂

![05.png](https://gw.alipayobjects.com/mdn/rms_5891a1/afts/img/A*u5rkS54zkgAAAAAAAAAAAAAAARQnAQ)

在查看streamFilter之前我们来看看工厂类是如何生产出限流器的

![06.png](https://gw.alipayobjects.com/mdn/rms_5891a1/afts/img/A*cj0nT5O69OYAAAAAAAAAAAAAARQnAQ)

限流器加入到限流链路结构中按照设定顺序依次生效。

CreateFilterChain 方法将多个filter 加入到链路结构中

![07.png](https://gw.alipayobjects.com/mdn/rms_5891a1/afts/img/A*a8ClQ76odpEAAAAAAAAAAAAAARQnAQ)

我们可以看到有各种各样的工厂类包括我们今天研究的限流工厂类实现了此接口

![08.png](https://gw.alipayobjects.com/mdn/rms_5891a1/afts/img/A*sBDbT44r2vgAAAAAAAAAAAAAARQnAQ)

### Stream_filter.go分析
![09.png](https://gw.alipayobjects.com/mdn/rms_5891a1/afts/img/A*wsw3RKe1GH8AAAAAAAAAAAAAARQnAQ)

## 整体流程：
最后我们再来回顾一下整体流程走向:

1. 从stream_filter_factory.go的初始化函数init() 开始,程序向creatorStreamFactory(map类型)插入了 createRpcFlowControlFilterFactory.

2. Mosn 创建出一个filter chain(代码位置[factory.go](https://github.com/mosn/mosn/tree/master/pkg/streamfilter/factory.go)) ,通过循环调用CreateFilterChain将所有的filter加入到链路结构包括我们今天的主人公限流器.

3. 创建限流器 NewStreamFilter().

4. 当流量通过mosn 将会进入到限流器的方法 OnReceive() 中并最终借助sentinel实现限流逻辑(是否已经达到阈值,是放行流量还是拦截流量, StreamFilterStop or StreamFilterContinue).


