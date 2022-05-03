# Layotto Commands Guidelines：

Layotto provides powerful commands, which makes contribution and local development/test easier. List as below:


## Highlights

+ You can simply run `make all` to format your codes, make style checks, linter checks, unit tests, and build layotto binary for host platform.

+ You can also run `make format` to format your codes. 

+ Run `make check` to make style checks, linter checks, unit tests.

+ Run `make build` to build layotto binary for host platform. 

See below commands to know more details or excute `make help`:

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