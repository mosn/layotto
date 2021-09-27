module mosn.io/layotto

go 1.14

require (
	github.com/dapr/components-contrib v1.2.0
	github.com/dapr/kit v0.0.1
	github.com/fsnotify/fsnotify v1.4.9
	github.com/gammazero/workerpool v1.1.2
	github.com/golang/mock v1.6.0
	github.com/golang/protobuf v1.5.2
	github.com/google/uuid v1.2.0
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0
	github.com/json-iterator/go v1.1.11
	github.com/pkg/errors v0.9.1
	github.com/shirou/gopsutil v3.21.3+incompatible
	github.com/stretchr/testify v1.7.0
	github.com/tetratelabs/proxy-wasm-go-sdk v0.1.2-0.20210520063156-d39281baed90
	github.com/tklauser/go-sysconf v0.3.5 // indirect
	github.com/urfave/cli v1.22.1
	github.com/valyala/fasthttp v1.26.0
	golang.org/x/sys v0.0.0-20210630005230-0f9fa26af87c // indirect
	google.golang.org/grpc v1.39.0
	google.golang.org/grpc/examples v0.0.0-20210917050107-e469f0d5f5bc // indirect
	google.golang.org/protobuf v1.27.1
	mosn.io/api v0.0.0-20210714065837-5b4c2d66e70c
	mosn.io/layotto/components v0.0.0-20210625065826-9c2ad8dbcf05
	mosn.io/layotto/sdk/go-sdk v0.0.0-20210926015636-a705edecedd9 // indirect
	mosn.io/layotto/spec v0.0.0-20210707123820-584778d048d3
	mosn.io/mosn v0.24.1-0.20210817133744-f3fc4c31ddee
	mosn.io/pkg v0.0.0-20210604065522-6e8f5a087814
	mosn.io/proxy-wasm-go-host v0.0.0-20210312032409-2334f9cf62ec
)

replace (
	mosn.io/layotto/components => ./components
	mosn.io/layotto/spec => ./spec
)
