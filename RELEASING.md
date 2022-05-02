# 发布手册
本文介绍下发布新版本时，发布负责人需要做什么

## 发布周期
Layotto 发布周期暂定为每季度发布一次。

## 发布 checklist

#### Step1: 检查当前迭代的 [Roadmap](https://github.com/mosn/layotto/projects) 
    
1. 检查进行中的任务
2. 检查未完成的任务
3. 与负责人确认任务状态和发布内容
  
#### Step2: 创建发布 tag， push 至 github 并检查工作流

1. 规范：请按照 `v{majorVersion}.{subVersion}.{latestVersion}` 格式创建 tag。
2. 等待 CI 结束，确认以下内容：
    + CI 测试 Jobs 全部通过：
        + Go 代码风格校验
        + Go 代码规范校验
        + Go 单元测试
        + Go 集成测试
    + CI 多平台 Artifacts 构建 Jobs 全部通过：
        + Linux/AMD64 Artifacts 成功 Build 并 Upload
        + Linux/ARM64 Artifacts 成功 Build 并 Upload
        + Darwin/AMD64 Artifacts 成功 Build 并 Upload
        + Darwin/ARM64 Artifacts 成功 Build 并 Upload
    + CI 多平台 Image 构建/发布 Jobs 全部通过：
        + Linux/AMD64 Image 成功 Build 并 Push DockerHub
        + Linux/ARM64 Image 成功 Build 并 Push DockerHub
            + Image Tag 规范：
                + AMD64/X86 架构的镜像：`layotto/layotto:{tag}`
                + ARM64 架构的镜像：`layotto/layotto.arm64:{tag}`

#### Step3: Draft a new release 并编写发布报告

> 发布报告可以先用 github 的功能自动生成，再基于生成的内容做修改。

> 可以参考以前的 [发版报告](https://github.com/mosn/layotto/releases)

#### Step4: 上传多平台架构的 Binaries

> 不必手动构建，直接将 `步骤 2` 中构建的多平台 Artifacts 下载上传即可

#### Step5: 确认发布

1. 点击发布
2. 社区周知
3. 检查 [Roadmap](https://github.com/mosn/layotto/projects)，修改上个版本未完成的任务，把 milestone 改为下个版本
4. 如果有 SDK 发布，需在 SDK 仓库做发布，并上传中央仓库 (比如 Java SDK 需要上传到 Maven 中央仓库)。

> TODO: need to translate.