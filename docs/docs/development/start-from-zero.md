# 新手攻略：从零开始成为 Layotto 贡献者
本文作为新手攻略，为想参与本开源项目、但是不清楚从哪下手的同学讲解下升级打怪路线。

## Step 1. Hello world,了解 Layotto 是干嘛的
先了解项目是干啥的，参照quickstart demo把hello world跑起来，比如：
   - [使用Configuration API](https://mosn.io/layotto/#/zh/start/configuration/start-apollo) 
   - [使用State API进行状态管理](https://mosn.io/layotto/#/zh/start/state/start) 
   - [使用分布式锁 API](https://mosn.io/layotto/#/zh/start/lock/start) 

Tips: 如果过程中遇到了报错没法启动，可以在github发issue提问。

Tips: 如果觉得某篇文档写的不够详细想要补充，或者文档写错了，可以提PR修复。这可能是成为contributor的最快路线 :)

（可选）扩展阅读：有一些介绍项目的演讲视频，比如[《MOSN 子项目 Layotto：开启服务网格+应用运行时新篇章》](https://mosn.io/layotto/#/zh/blog/mosn-subproject-layotto-opening-a-new-chapter-in-service-grid-application-runtime/index) ， 
比如 [《Service Mesh落地之后：为 sidecar 注入灵魂》](https://www.bilibili.com/video/BV1RL4y1b7U9?from=search&seid=1492521025214444985&spm_id_from=333.337.0.0)
。不过看视频比较花时间，懒得看的话可以多跑几个demo体验下，quickstart demo蛮多的

## Step 2. 挑选适合你的任务
当你大概知道这个项目是做什么的了，可以从 [good first issue 列表](https://github.com/mosn/layotto/issues?q=is%3Aissue+is%3Aopen+sort%3Aupdated-desc+label%3A%22good+first+issue%22) 中挑选感兴趣的任务，没打勾的任务都可以认领。完成任务即可成为Layotto Contributor。
选择任务时，推荐：
   
- 对于想要学习go语言的同学 

可以先选择"给指定模块添加注释"的任务，读代码的同时加注释，顺便学习go写法；

或者选择加单测的任务，加单测的过程可以练习写go；

同时，Easy级别会有一些写代码的任务，写起来并不难但是对项目很有帮助，感兴趣可以直接上手写代码，比如 https://github.com/mosn/layotto/issues/275#issuecomment-957711746 比如 [开发 in-memory 组件](https://github.com/mosn/layotto/issues/67#issuecomment-975134341) 具体可以留意任务列表中的Easy任务

- 对于想要写其他语言的同学（java,c++,python,typescript等）

虽然Layotto是 go语言开发的项目，但是需要开发各种语言的sdk，而开发sdk是不需要懂go语言的。
比如熟悉java的同学可以认领java sdk相关的任务，并不需要懂go。

因此可以留意任务列表中和sdk相关的任务。


- 对于有一定后端基础、go语言基础的同学

有基础的话，按兴趣选择任务即可，比如对分布式锁感兴趣，比如对Webassembly感兴趣，可以选择跟这些技术相关的任务。

## Step 3. 认领任务之后……沟通很重要！
认领任务后，如果是有一定开发成本的任务，最好先在issue下面描述下自己的设计方案，避免返工。

开发过程中难免遇到困难，比如有报错、一直修不好，这都很正常，因此懂得向他人求助、求助的时候描述清楚自己的问题也是很重要的。

遇到问题可以在issue区或者钉钉群里讨论，可以描述下自己遇到的异常现象（比如报错信息）、复现步骤，以便大家帮忙排查。
