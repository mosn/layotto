module github.com/layotto/layotto

go 1.14

require (
	github.com/dapr/components-contrib v1.1.0-rc1
	github.com/dapr/dapr v1.0.1-0.20210325161510-849f52560a63
	github.com/golang/mock v1.4.4
	github.com/golang/protobuf v1.5.2
	github.com/google/uuid v1.2.0
	github.com/json-iterator/go v1.1.10
	github.com/pkg/errors v0.9.1
	github.com/shirou/gopsutil v3.21.3+incompatible
	github.com/sirupsen/logrus v1.6.0
	github.com/stretchr/testify v1.7.0
	github.com/tklauser/go-sysconf v0.3.5 // indirect
	github.com/urfave/cli v1.22.1
	github.com/valyala/fasthttp v1.21.0
	github.com/zouyx/agollo/v4 v4.0.6
	go.etcd.io/etcd v3.3.25+incompatible
	google.golang.org/grpc v1.37.0
	google.golang.org/protobuf v1.26.0
	mosn.io/api v0.0.0-20210414070543-8a0686b03540
	mosn.io/mosn v0.22.1-0.20210425073346-b6880db4669c
	mosn.io/pkg v0.0.0-20210401090620-f0e0d1a3efce
)

replace (
	go.etcd.io/etcd => go.etcd.io/etcd v0.0.0-20210226220824-aa7126864d82
	google.golang.org/genproto => google.golang.org/genproto v0.0.0-20200305110556-506484158171
	google.golang.org/grpc => google.golang.org/grpc v1.28.0
	k8s.io/client => github.com/kubernetes-client/go v0.0.0-20190928040339-c757968c4c36
)
