# 发布手册
本文介绍下发布新版本时，发布负责人需要做什么
## 发布周期
Layotto 发布周期暂定为每季度发布一次。

## 发布 checklist

- [ ] 检查 [当前迭代的roadmap](https://github.com/mosn/layotto/projects) ，看看是否有进行中、没做完的事情，和负责人确认下能否带着一起发布  
  
- [ ] 跑通所有测试（包括单元测试，集成测试和demo）
  
- [ ] 编译不同操作系统下的二进制 Layotto。 至少包括 Linux和 Mac(amd64)的:

![img.png](../../img/development/release/img.png)
  
- [ ] Draft a new release

![img_1.png](../../img/development/release/img_1.png)
  
- [ ] 打tag、写发布报告

发布报告可以先用github的功能自动生成，再基于生成的内容做修改。

可以参考以前的 [发版报告](https://github.com/mosn/layotto/releases)

- [ ] 点按钮，发布
- [ ] 群里发个消息通知下大家  
- [ ] 发布完成后，创建[下次迭代的roadmap](https://github.com/mosn/layotto/projects) （可能已经有了），可以把上个版本没做完的事项挪进来
- [ ] 如果有相应的 SDK 需要发布，参考上述流程，在 SDK 仓库做发布，并上传中央仓库 (比如 java sdk需要上传到 Maven 中央仓库)。