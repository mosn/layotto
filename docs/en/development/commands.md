# Layotto Commands Guidelines：

Layotto provides powerful commands, which makes contribution and local development/test easier. List as below:


## Highlights

+ You can simply run `make all` to format your codes, make style checks, linter checks, unit tests, and build layotto binary for host platform.

+ You can also run `make format` to format your codes. 

+ Run `make check` to make style checks, linter checks, unit tests.

+ Run `make build` to build layotto binary for host platform. 

+ Run `make license` to add licnese headers for all code files, with docker.

> Notice: If you encounter the error "make[1]: *** No rule to make target 'all'. Stop.", it means that the Makefile cannot find this target.
>
> You should refer to the commands provided by `make help`. For example,  You can execute commands like `make lint` `make format` `make test` in place of `make all`, which you think should be used to checked codes locally.

See below commands to know more details or excute `make help`:

```
Layotto is an open source project for a fast and efficient cloud native application runtime.

Usage:
  make <Target> <Option>

Targets:

Golang Development
  build            Build layotto for host platform.
  multiarch        Build layotto for multiple platforms.
  clean            clean all unused generated files.
  lint             Run go syntax and styling of go sources.
  test             Run golang unit test in target paths.
  workspace        check if workspace is clean and committed.
  format           Format codes style with gofmt and goimports.

Image Development
  image            Build docker images for host arch.
  image-multiarch  Build docker images for multiple platforms.
  push             Push docker images to registry.
  push-multiarch   Push docker images for multiple platforms to registry.
  proxyv2          Build proxy image for host arch.
  proxyv2-push     Push proxy image to registry.

Proto Development
  proto            Generate code and documentation based on the proto files.
  proto-doc        Generate documentation based on the proto files.
  proto-code       Generate code based on the proto files.
  proto-lint       Run Protobuffer Linter with Buf Tool

Kubernetes Development
  deploy-k8s       Install Layotto in Kubernetes.
  undeploy-k8s     Uninstall Layotto in Kubernetes.

WebAssembly Development
  wasm-build       Build layotto wasm for linux arm64 platform.
  wasm-image       Build layotto wasm image for multiple platform.
  wasm-push        Push layotto wasm image for multiple platform.

CI/CD Development
  base             Build base docker images for host arch.
  base-multiarch   Build base docker images for multiple platforms.
  deadlink         Run deadlink check test.
  quickstart       Run quickstart check test.
  coverage         Run coverage analysis.
  license          Add license headers for code files.
  license-check    Check codes license headers.
  integrate-wasm   Run integration test with wasm.
  integrate-runtime  Run integration test with runtime.


Options:

  BINS         The binaries to build. Default is all of cmd.
               This option is available when using: make build/multiarch
               Examples:
               * make multiarch BINS="layotto"
               * make build BINS="layotto_multiple_api layotto"
  IMAGES       Backend images to make. Default is all of cmds.
               This option is available when using: make image/image-multiarch/push/push-multiarch
               Examples: 
               * make image IMAGES="layotto"
               * make image-multiarch IMAGES="layotto"
               * make push IMAGES="layotto_multiple_api"
               * make push-multiarch IMAGES="layotto_multiple_api"
  NAMESPACE    The namepace to deploy. Default is `default`.
               This option is available when using: make deploy-k8s/undeploy-k8s
               Examples: 
               * make deploy-k8s NAMESPACE="layotto"
               * make undeploy-k8s NAMESPACE="default"
  VERSION    The image tag version to build. Default is the latest release tag.
               This option is available when using: make image/image-multiarch/push/push-multiarch
               Examples: 
               * make image VERSION="latest"
               * make image-multiarch VERSION="v1.0.0"
               * make push-multiarch VERSION="v2.0.0"
  REGISTRY_PREFIX    The docker image registry repo name to push. Default is `layotto`.
               This option is available when using: make push/push-multiarch
               Examples: 
               * make push IMAGES="layotto" REGISTRY_PREFIX="mosn"
               * make push IMAGES="layotto_multiple_api" REGISTRY_PREFIX="mosn"
               Supported Platforms: linux_amd64 linux_arm64 darwin_amd64 darwin_arm64
  PLATFORMS    The multiple platforms to build. Default is linux_amd64 and linux_arm64.
               This option is available when using: make multiarch
               Examples: 
               * make multiarch BINS="layotto" PLATFORMS="linux_amd64 linux_arm64"
               Supported Platforms: linux_amd64 linux_arm64 darwin_amd64 darwin_arm64
```