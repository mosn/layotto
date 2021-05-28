module github.com/layotto/layotto

go 1.14

require (
	github.com/golang/mock v1.4.4
	github.com/golang/protobuf v1.5.2
	github.com/pkg/errors v0.9.1
	github.com/shirou/gopsutil v3.21.3+incompatible
	github.com/stretchr/testify v1.7.0
	github.com/tetratelabs/proxy-wasm-go-sdk v0.1.2-0.20210520063156-d39281baed90
	github.com/tklauser/go-sysconf v0.3.5 // indirect
	github.com/urfave/cli v1.20.0
	github.com/zouyx/agollo/v4 v4.0.6
	go.etcd.io/etcd v3.3.25+incompatible
	google.golang.org/grpc v1.37.0
	google.golang.org/protobuf v1.26.0
	mosn.io/api v0.0.0-20210414070543-8a0686b03540
	mosn.io/mosn v0.22.1-0.20210425073346-b6880db4669c
	mosn.io/pkg v0.0.0-20210401090620-f0e0d1a3efce
	mosn.io/proxy-wasm-go-host v0.0.0-20210312032409-2334f9cf62ec
)

replace (
	github.com/golang/protobuf => github.com/golang/protobuf v1.3.5
	go.etcd.io/etcd => go.etcd.io/etcd v0.0.0-20210226220824-aa7126864d82
	google.golang.org/grpc => google.golang.org/grpc v1.28.0
)
