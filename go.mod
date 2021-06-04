module github.com/layotto/layotto

go 1.14

require (
	github.com/dapr/components-contrib v1.2.0
	github.com/dapr/kit v0.0.1
	github.com/golang/mock v1.4.4
	github.com/google/uuid v1.2.0
	github.com/json-iterator/go v1.1.10
	github.com/layotto/layotto/components v0.0.0-20210603045430-66065fa0b67f
	github.com/layotto/layotto/spec v0.0.0-20210604023314-bb30491493a4
	github.com/shirou/gopsutil v3.21.3+incompatible
	github.com/stretchr/testify v1.7.0
	github.com/tklauser/go-sysconf v0.3.5 // indirect
	github.com/urfave/cli v1.22.1
	github.com/valyala/fasthttp v1.26.0
	google.golang.org/grpc v1.37.0
	google.golang.org/grpc/examples v0.0.0-20210526223527-2de42fcbbce3 // indirect
	google.golang.org/protobuf v1.26.0
	mosn.io/api v0.0.0-20210414070543-8a0686b03540
	mosn.io/mosn v0.22.1-0.20210425073346-b6880db4669c
	mosn.io/pkg v0.0.0-20210401090620-f0e0d1a3efce
)

replace (
	github.com/layotto/layotto/components v0.0.0-20210603045430-66065fa0b67f => ./components
	github.com/layotto/layotto/spec v0.0.0-20210604023314-bb30491493a4 => ./spec
)
