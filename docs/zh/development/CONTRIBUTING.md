# Layotto 贡献指南

Layotto 基于 Apache 2.0 许可发布，遵循标准的 Github 开发流程，使用 Github Issue 跟踪问题并将 Pull Request 合并到 main 分支中。本文档可以帮助你了解如何参与贡献。

## 贡献流程
- 首次提交 Pull Request 后，您需要先签署[贡献者许可协议（CLA）](http://cla.sofastack.tech/mosn)
- 优化您的 Pull Request，确保能通过自动化测试(CI) 
  - 如果是首次贡献者，您提交 Pull Request 是没法自动触发 CI 的，需要由项目维护者手动运行 CI. 这是 Github 的默认限制, 但您做过一次贡献、成为 contributor 后，再提新 PR 就能自动触发 CI 了。
  - CI 的详细说明见 [Layotto GitHub Workflows](zh/development/github-workflows)
  - 为了方便开发， Layotto 有一套 make 脚本，可以在本地跑检查、跑测试，在本地启动好 docker 后，敲 `make all` 即可，详见[文档说明](zh/development/commands)
- 由社区维护者来 code review
  - code review 如果有修改意见，会在 PR 下留言
  - code review 有两票 approve 后, PR 就可以合并

## 代码约定

以下对于 Pull Request 的要求并非强制，但是会对您提交 Pull Request 有帮助。

1. 代码格式规范: 你只需要执行 `make format` 去格式化你的代码。
2. 确保所有新的 `.go` 文件都具有简单的 doc 类注释
3. 将 Apache Software Foundation 许可证标头注释添加到所有新的 `.go` 文件（可以从项目中的现有文件复制）
4. 添加文档
5. 进行一些单元测试也会有很大帮助。
6. 请确保代码覆盖率不会降低。
7. 确保提交 Pull Request  之前所有的测试能够正确通过。你可以本地启动好 docker 后，执行 `make all` 去格式化你的代码，进行风格测试，linter 规范测试，单元测试，以及集成测试。但是执行前[请先查看注意事项](zh/development/commands)
9. 按照 Github 工作流提交 Pull Request  ，并遵循 Pull Request 的规则。

> Layotto 提供了很多方便本地开发的命令行工具，请在[这里](zh/development/commands)进行查阅

## 版本命名约定

Layotto 的版本包含三位数，格式为 x.x.x，第一位是出于兼容性考虑； 第二个是新功能和增强功能； 最后一位是错误修复。

## 维护者代码审查策略

项目维护者审查代码时建议遵循以下策略：

1. 检查 PR 对应的 Issue
2. 检查解决方案的合理性
3. 检查 UT 和 Benchmark 的结果
4. 注意使代码结构发生变化的代码，全局变量的用法，特殊情况和并发的处理