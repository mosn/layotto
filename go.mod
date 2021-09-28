module mosn.io/layotto

go 1.14

require (
	github.com/dapr/components-contrib v1.4.0-rc2
	github.com/dapr/kit v0.0.2-0.20210614175626-b9074b64d233
	github.com/fsnotify/fsnotify v1.4.9
	github.com/gammazero/workerpool v1.1.2
	github.com/golang/mock v1.6.0
	github.com/golang/protobuf v1.5.2
	github.com/google/cel-go v0.7.3 // indirect
	github.com/google/uuid v1.2.0
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/json-iterator/go v1.1.11
	github.com/pkg/errors v0.9.1
	github.com/shirou/gopsutil v3.21.3+incompatible
	github.com/stretchr/testify v1.7.0
	github.com/urfave/cli v1.22.1
	github.com/valyala/fasthttp v1.28.0
	google.golang.org/grpc v1.39.0
	google.golang.org/grpc/examples v0.0.0-20210917050107-e469f0d5f5bc // indirect
	google.golang.org/protobuf v1.27.1
	mosn.io/api v0.0.0-20210714065837-5b4c2d66e70c
	mosn.io/layotto/components v0.0.0-20210625065826-9c2ad8dbcf05
	mosn.io/layotto/spec v0.0.0-20210707123820-584778d048d3
	mosn.io/mosn v0.24.1-0.20210928035557-38b3b922b595
	mosn.io/pkg v0.0.0-20210823090748-f639c3a0eb36
	mosn.io/proxy-wasm-go-host v0.0.0-20210312032409-2334f9cf62ec
)

replace (
	github.com/google/cel-go => github.com/google/cel-go v0.5.1
	github.com/tetratelabs/proxy-wasm-go-sdk => github.com/layotto/proxy-wasm-go-sdk v0.14.1-0.20210926122819-378cc27b0ffb
	mosn.io/layotto/components => ./components
	mosn.io/layotto/spec => ./spec
)
