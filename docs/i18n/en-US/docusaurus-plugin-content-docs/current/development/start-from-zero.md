# Newhanded offensive：became a Layotto contributor from zero

This paper is a new hand-walker for the purpose of attempting to participate in this open source project, but it is not clear from what classmates know about which to upgrade.

## Step 1. Hello world, understand Layotto is dry

Learn what the project was. Use quickstart to run hello around the world by reference to quickstart demo like：

- [Use configuration API](https://mosn.io/layotto/en-US/docs/start/configuration/start-apollo)
- [State API for status management](https://mosn.io/layotto/#/start/state/start)
- [Use distributed lock API](https://mosn.io/layotto/#/start/lock/start)

Tips: Ask issues in github, if a bug is started during the process.

Tips: If you feel that a document is not written in sufficient detail to add it, or the document is written incorrectly, you can make a PR.This is probably the fastest route to contribor :)

（可选）扩展阅读：有一些介绍项目的演讲视频，比如[《MOSN 子项目 Layotto：开启服务网格+应用运行时新篇章》](https://mosn.io/layotto/docs/blog/mosn-subproject-layotto-opening-a-new-chapter-in-service-grid-application-runtime/index) ，
比如 [《Service Mesh落地之后：为 sidecar 注入灵魂》](https://www.bilibili.com/video/BV1RL4y1b7U9?from=search\&seid=1492521025214444985\&spm_id_from=333.337.0.0)
。But looking at videos takes time, lazy to run more than a few demo experiments, quickstart more brutally

## Step 2. Select your task

当你大概知道这个项目是做什么的了，可以从 [good first issue 列表](https://github.com/mosn/layotto/issues?q=is%3Aissue+is%3Aopen+sort%3Aupdated-desc+label%3A%22good+first+issue%22) 中挑选感兴趣的任务，没打勾的任务都可以认领。Complete task to become Layotto Contribor.
Recommended： when selecting a task

- Recourse to the language you want to learn

You can first select the task "Add a comment to the specified module", read the code with annotations, and learn the towriting method.

Or select a single task with a single metering process that can be used to write go;

At the same time, the Easy level will have some code writing tasks that are not difficult to write but are very helpful to the project. Interest can be taken directly to the handwritten code, such as https://github.com/mosn/layotto/issues/275#issuomment-957711746 like [Developing in-memory components] (https://github.com/mosn/layotto/issues/67#issuomment_975134341)

- For students who want to write other languages (java, c++, python, typescript, etc.)

Although Layotto is a project for development of the language of go, there is a need to develop sdk in various languages, which does not need to understand the language of go.
For example, students familiar with java can claim java sdk and do not need to know go.

It is therefore possible to look at tasks in the task list related to sdk.

- Recourse to some back-end base, go-language base

If available, select tasks by interest, such as interested in distributive locks, such as Webassembly, and select tasks related to these technologies.

## Step 3. 认领任务之后……沟通很重要！

When it is assumed that there is a certain development cost, it is best to describe its own design programme below the issue and avoid returning to work.

The difficulties encountered in the development process, such as misreporting and continuing misconduct, were normal, and it was therefore important to know when asking and asking for help from others to describe their own problems.

Problems can be discussed in an issue area or pegged group, which can describe the anomalies that they encounter (e.g. misinformation for misinformation) and the steps taken to replicate them so that they can be easily checked.
