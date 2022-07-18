# SOFAStack & MOSN 社区角色说明 v2.0
![](https://gw.alipayobjects.com/mdn/rms_95b965/afts/img/A*DpjGQqAcRyQAAAAAAAAAAAAAARQnAQ)

Layotto是MOSN社区的一个子项目，本文介绍SOFAStack & MOSN社区角色的 v2.0提案

该提案会先在Layotto和SOFA-Registry项目试行

## 0. 设计目标
- 吸引更多的人加入成为社区的**维护者、管理者**，而不仅仅是由核心开发者维护社区、单点瓶颈
- 提升社区活跃度

## 1. 原则
### 1.0. 来了就是朋友
降低门槛，让感兴趣的爱好者能参与到社区维护、管理中

### 1.1. Title 社区化：所有社区 title、权限和参与者所属公司/组织无关
为了调动大家积极性，让项目 owner 觉得是为了个人爱好/个人荣誉而参与社区，不管去哪家公司都还是SOFAStack PMC member/Co-founder; 让每个感兴趣的技术同学都能参与到社区核心决策中，促进社区活跃。

比如某位开发者创建了项目A，多年后即使他不再参与社区维护、辗转跳槽多家公司，经历了风风雨雨、起起落落，决定转行不做程序员了、开了家火锅店，他也依然是项目A的 Co-founder。

比如某位编程爱好者主业是火锅店老板、不是程序员，只要他对项目的贡献足够，他也可以是PMC member 或项目 Reviewer。

### 1.2. 确定性回报：制定明确的社区成员晋级规则
社区成员做出怎样的贡献后可以晋级？要有透明、具体的规则，让社区同学从 Contributor 成长到 Committer 或 Maintainer 有路可循，**让大家觉得“投入就一定有回报”**
每个项目制定自己的晋升规则，但要明确做了哪些事情就可以晋升，让大家看到投入社区的回报

### 1.3. 决策社区化
#### 评审社区化：需求和技术方案公开评审
重要需求的评审、技术方案的评审应该在社区内**以开放的方式做评审，而不是几个人搞内部评审会**

例如要做技术方案评审时，社区成员可以在自己的用户群或者 Core Team 群里组织线上会议、开直播，允许每一个感兴趣的人观看、参与。如果觉得直播太重，issue 区发 proposal 也行。

反例：张三是一位火锅店老板，同时也是某项目的Committer。他很有参与社区建设的热情，但是技术方案评审都不带他，他感到很失望，这个Committer title形同虚设。

#### 晋升社区化
晋升规则、提名由社区成员在群内协商决定


## 2. 新角色解释
### 2.0. Co-founder（项目创始人）终身荣誉制
#### 角色说明
作为终身荣誉，但没有实权。
比如某位开发者创建了项目A，多年后即使他不再参与社区维护、转行不做程序员了，他也依然是项目A的 Co-founder。

这是一种荣誉和激励，希望让大家觉得是在为自己的个人影响力做开源
#### 设计这个角色的目的
激励大家往SOFAStack和MOSN社区开源新项目（下半年已经有几个同学有想法了，而且sofa-bolt-go已经在开发了 )；

激励当前的项目owner（如果是创始人的话）；

感谢已经离开社区的初始开发者，邀请回到社区

### 2.1. Member
#### 角色说明
新增Member角色，其实就是邀请加入组织，加入后个人档案页会展示logo。

当然，考虑到过往一些项目是只有Committer才能拉组织、门槛比较高，该角色可选、不强制要求，各项目可以按自己项目的情况决定是否加入该角色。
目前的想法是先在Layotto和SOFA-Registry试验加入该角色。
#### 添加这个角色的目的?
加入组织后，贡献者的github profile页显示组织logo，这很酷，可以激励大家加入社区、做出贡献。
#### Member的职责
Member职责是code review以及处理issue，如果member处理issue和code review的频率很高，就可以成为Reviewer（详见下）；

#### Member分配的github组织权限
默认权限，没有写权限。但是LGTM有效
#### 成为Member的条件（加入组织的条件）
每个项目自己定
### 2.2. Reviewer
#### 角色说明
和Member一样，也是负责code review以及处理issue，算是对Member贡献的认可和感谢。长期review贡献多可以晋升Committer

Reviewer要承担一定的review责任。比如某个issue/Pull request涉及某个模块，项目维护者可以request这个模块的Reviewer进行Review、负责把关。比如有人给Layotto的WASM模块提bug，可以request WASM模块的reviewer @zhenjunMa

#### 成为Reviewer的条件
每个项目自己定
#### Reviewer是否会回退到Member
每个项目自己定。目前收集到的不同意见：

a. 成为Reviewer之后，职责相比member要更明确，就是处理issue和PR，而且必须保持每周至少处理2个issue或1个PR，如果一个月内，有两周处理的频率没有达到标准，会回退到Member；

b. 认为给Reviewer定目标太有KPI味了，工作背KPI就够难了别给社区搞KPI
#### Reviewer分配的github组织权限
![](https://user-images.githubusercontent.com/26001097/129857585-2f2ddcda-4a5d-4f94-a36d-48e6a9f52e0e.png)

#### 添加这个角色的目的?
1. 希望能让MOSN和SOFAStack社区响应速度更快一些，响应pr快一些（至于如何提高issue响应速度，我还没想到啥办法，欢迎提建议

2. 决策社区化，code review可以由活跃贡献者来，不想只有一两个人有review权限
   
3. 很多爱好者可能一段时间有空、一段时间没空（国内程序员现状.....)，没空的时候通过Review PR和issue，可以只花少量时间就能参与项目、做出贡献、知道项目这段时间发生的变化

#### code review时，一票LGTM还是两票LGTM才能合并？
A: 目前想法是参考vue社区，如果review通过可以回复LGTM，但是没有merge权限，项目维护者（有写权限的人）负责实际合并。
项目维护者如果很信任reviewer，可以简单看下pr，没问题就合并；如果pr改动比较大、不太放心，就要自己也认真看一下pr。

所以，约等于两票合并制。

这么设计是想迫使项目核心维护者了解每个pr、每个改动。

## 3. 其他
尽量在issue和PR使用英语

用哪种IM群由每个项目自行讨论决定，可选的包括微信群，Slack,gitter

社区会议每个项目自己定

## 4. 参考资料
Apollo
[https://github.com/apolloconfig/apollo/pull/3670](https://github.com/apolloconfig/apollo/pull/3670)
[https://github.com/apolloconfig/apollo/issues/3684](https://github.com/apolloconfig/apollo/issues/3684)
[https://github.com/apolloconfig/apollo/discussions/categories/announcements](https://github.com/apolloconfig/apollo/discussions/categories/announcements)


Tidb的社区组织架构 [https://pingcap.com/blog-cn/tidb-community-upgrade/](https://pingcap.com/blog-cn/tidb-community-upgrade/)