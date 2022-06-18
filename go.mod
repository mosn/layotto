module mosn.io/layotto

go 1.14

require (
	github.com/99designs/keyring v1.2.1 // indirect
	github.com/Azure/go-autorest/autorest/azure/cli v0.4.2 // indirect
	github.com/Azure/go-autorest/autorest/validation v0.3.1 // indirect
	github.com/Azure/go-autorest/logger v0.2.1 // indirect
	github.com/SkyAPM/go2sky v0.5.0
	github.com/agrea/ptr v0.0.0-20180711073057-77a518d99b7b
	github.com/alicebob/miniredis/v2 v2.16.0
	github.com/dapr/components-contrib v1.5.1-rc.1
	github.com/dapr/kit v0.0.2-0.20210614175626-b9074b64d233
	github.com/dimchansky/utfbom v1.1.1 // indirect
	github.com/fsnotify/fsnotify v1.4.9
	github.com/gammazero/workerpool v1.1.2
	github.com/golang/mock v1.6.0
	github.com/golang/protobuf v1.5.2
	github.com/google/uuid v1.3.0
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/json-iterator/go v1.1.12
	github.com/opentracing/opentracing-go v1.2.0
	github.com/openzipkin/zipkin-go v0.4.0
	github.com/phayes/freeport v0.0.0-20180830031419-95f893ade6f2
	github.com/pkg/errors v0.9.1
	github.com/shirou/gopsutil v3.21.3+incompatible
	github.com/stretchr/testify v1.7.0
	github.com/uber/jaeger-client-go v2.25.0+incompatible
	github.com/urfave/cli v1.22.1
	github.com/valyala/fasthttp v1.28.0
	go.uber.org/automaxprocs v1.4.0 // indirect
	golang.org/x/tools v0.1.10 // indirect
	google.golang.org/grpc v1.39.0
	google.golang.org/grpc/examples v0.0.0-20210818220435-8ab16ef276a3
	google.golang.org/protobuf v1.27.1
	mosn.io/api v1.0.0
	mosn.io/layotto/components v0.0.0-20220413092851-55c58dbb1a23
	mosn.io/layotto/spec v0.0.0-20220413092851-55c58dbb1a23
	mosn.io/mosn v1.0.2-0.20220607053200-448da0ebe355
	mosn.io/pkg v1.0.0
	mosn.io/proxy-wasm-go-host v0.1.1-0.20210524020952-3fb13ba763a6
	nhooyr.io/websocket v1.8.7 // indirect
)

replace (
	github.com/gin-gonic/gin => github.com/gin-gonic/gin v1.7.0
	github.com/klauspost/compress => github.com/klauspost/compress v1.13.0
	mosn.io/layotto/components => ./components
	mosn.io/layotto/spec => ./spec
	mosn.io/proxy-wasm-go-host => github.com/layotto/proxy-wasm-go-host v0.1.1-0.20210929091514-828451606147
)
