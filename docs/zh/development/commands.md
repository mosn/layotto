# Layotto 命令指南：

Layotto 提供了丰富的命令行工具，方便贡献者开发和测试 Layotto 代码，具体命令如下：


## 重要开发命令

+ 你可以执行 `make all` 去格式化你的代码，进行风格测试，linter 规范测试，单元测试，以及构建当前平台的二进制文件。

+ 你也可以执行 `make format` 去格式化你的代码

+ 执行 `make check` 进行风格测试，linter 规范测试，单元测试

+ 执行 `make build` 构建当前平台的二进制文件

具体细节可查看一下命令，或执行 `make help` 查看：

```
Layotto commands 👀: 

A fast and efficient cloud native application runtime 🚀.
Commands below are used in Development 💻 and GitHub workflow 🌊.

Usage: make <COMMANDS> <ARGS> ...

COMMANDS:
  build               Build layotto for host platform.
  build.multiarch     Build layotto for multiple platforms. See option PLATFORMS.
  image               Build docker images for host arch.
  image.multiarch     Build docker images for multiple platforms. See option PLATFORMS.
  push                Push docker images to registry.
  push.multiarch      Push docker images for multiple platforms to registry.
  app                 Build app docker images for host arch. [`/docker/app` contains apps dockerfiles]
  app.multiarch       Build app docker images for multiple platforms. See option PLATFORMS.
  wasm                Build layotto wasm for linux arm64 platform.
  wasm.multiarch      Build layotto wasm for multiple platform.
  wasm.image          Build layotto wasm image for multiple platform.
  wasm.image.push     Push layotto wasm image for multiple platform.
  check               Run all go checks of code sources.
  check.lint          Run go syntax and styling of go sources.
  check.unit          Run go unit test.
  check.style         Run go style test.
  style.coverage      Run coverage analysis.
  style.deadlink      Run deadlink check test.
  style.quickstart    Run quickstart check test.
  integrate.wasm      Run integration test with wasm.
  integrate.runtime   Run integration test with runtime.
  format              Format layotto go codes style with gofmt and goimports.
  clean               Remove all files that are created by building.
  all                 Run format codes, check codes, build Layotto codes for host platform with one command
  help                Show this help info.

ARGS:
  BINS         The binaries to build. Default is all of cmd.
               This option is available when using: make build/build.multiarch
               Example: make build BINS="layotto_multiple_api layotto"
  IMAGES       Backend images to make. Default is all of cmds.
               This option is available when using: make image/image.multiarch/push/push.multiarch
               Example: make image.multiarch IMAGES="layotto_multiple_api layotto"
  PLATFORMS    The multiple platforms to build. Default is linux_amd64 and linux_arm64.
               This option is available when using: make build.multiarch/image.multiarch/push.multiarch
               Example: make image.multiarch IMAGES="layotto" PLATFORMS="linux_amd64 linux_arm64"
               Supported Platforms: linux_amd64 linux_arm64 darwin_amd64 darwin_arm64
```