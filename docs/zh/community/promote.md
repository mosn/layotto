# Layotto社区晋升规则
![](https://gw.alipayobjects.com/mdn/rms_95b965/afts/img/A*DpjGQqAcRyQAAAAAAAAAAAAAARQnAQ)

## 1. Member
### 成为Member的条件
满足以下条件可以申请成为Member:
- 贡献过一个有价值的PR，例如一个 Easy 级别的社区开发任务
- 有意愿一起维护社区

### 如何申请成为 Member
可以在 [discussion 区](https://github.com/mosn/layotto/discussions) 发个帖、简单自我介绍下，比如提交过哪些pr, 对哪方面的技术感兴趣，后续有兴趣一起做哪方面的贡献。如果不想透露隐私，就不写个人信息。

这个环节的目的是让大家互相认识一下，简单写几句话即可。

参考 https://github.com/mosn/layotto/discussions/675

### 职责
Member 需要一起帮忙回复issue/pr，triage（把issue分配给对应模块的负责人）

### 权限
Triage权限。有权限操作issue和pr，例如打label、分配问题。
详细的权限说明见 [permissions-for-each-role](https://docs.github.com/en/organizations/managing-access-to-your-organizations-repositories/repository-roles-for-an-organization#permissions-for-each-role)

## 2. Reviewer
### 成为Reviewer的条件
有意愿负责某个模块的issue review和code review，且对该模块贡献过的PR满足**下列条件之一:**
- 1个Hard级别的PR
- 2个Medium级别的PR
- 1个Medium+2个Easy级别的PR

注：相当于`Hard:Medium:Easy`的换算关系是`1:2:4`

设计这个规则的逻辑是： Reviewer要对某个模块很懂，才能对这个模块把关。那怎么判断他很懂呢？就是看他做过的PR，1个hard级别的pr，或者2个medium级别的pr，或者1个medium+2个easy级别的pr

这么设计的缺点：没有把大家做review的贡献纳入进来，没有激励参与者做review。也考虑了把review次数纳入晋升条件的话，但是仔细想想不太好统计，这方面大家有啥建议欢迎讨论。

### 职责
负责某个模块的issue review和code review,给出技术建议。有该模块相关的重大变更会request review模块Reviewer。

## 3. Committer
### 成为Committer的条件
贡献过的PR满足下列条件:
- 合并的 PR 达到 10个；
- 其中至少包含1个 Hard 级别PR, 或者4个 Medium 级别PR；

### 职责
- 社区咨询支持；
- 积极响应 assign 给您的 Issue 或 PR；
- 对于社区重大决定的投票权；
- Review 社区的 PR；

### 认证、运营宣传
- 在Discussion区颁发电子证书

示例： [Welcome new committer: Zhang Li Bin](https://github.com/mosn/layotto/discussions/352)
  
- 邮寄实体证书

- 公众号宣传

示例：

[恭喜 张立斌 成为 Layotto committer！](https://mp.weixin.qq.com/s/no6mDymNEGxH3uoZbl1YTQ)

[恭喜 赵延 成为 SOFAJRaft committer！](https://mp.weixin.qq.com/s/BKJ0bcaGBeYNErDhpjk42Q)

## 4. PMC
项目管理委员会，为项目核心管理团队，参与 roadmap 制定与社区相关的重大决议；

### 加入 PMC 的条件
由PMC Member为某位Committer提名，然后PMC 投票，投票过半即可晋升为PMC Member

### 职责
积极参与社区讨论，对社区重大决策给予指导；
负责保证开源项目的社区活动都能运转良好；