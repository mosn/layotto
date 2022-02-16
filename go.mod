module mosn.io/layotto

go 1.14

require (
	github.com/Azure/go-autorest/autorest/azure/cli v0.4.2 // indirect
	github.com/Azure/go-autorest/autorest/validation v0.3.1 // indirect
	github.com/Azure/go-autorest/logger v0.2.1 // indirect
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
	github.com/json-iterator/go v1.1.11
	github.com/phayes/freeport v0.0.0-20180830031419-95f893ade6f2
	github.com/pkg/errors v0.9.1
	github.com/shirou/gopsutil v3.21.3+incompatible
	github.com/stretchr/testify v1.7.0
	github.com/urfave/cli v1.22.1
	github.com/valyala/fasthttp v1.28.0
	go.uber.org/automaxprocs v1.4.0 // indirect
	golang.org/x/crypto v0.0.0-20210921155107-089bfa567519 // indirect
	golang.org/x/net v0.0.0-20211005001312-d4b1ae081e3b // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/grpc v1.39.0
	google.golang.org/grpc/examples v0.0.0-20210818220435-8ab16ef276a3
	google.golang.org/protobuf v1.27.1
	mosn.io/api v0.0.0-20211217011300-b851d129be01
	mosn.io/layotto/components v0.0.0-20220119065745-4f03f6779399
	mosn.io/layotto/spec v0.0.0-20220119065745-4f03f6779399
	mosn.io/mosn v0.25.1-0.20211217125944-69b50c40af81
	mosn.io/pkg v0.0.0-20211217101631-d914102d1baf
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
