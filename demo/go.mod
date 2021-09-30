module mosn.io/layotto/demo

go 1.14

require (
	github.com/golang/protobuf v1.5.0
	github.com/google/uuid v1.2.0
	github.com/tetratelabs/proxy-wasm-go-sdk v0.14.1-0.20210922004205-46e3ac3a25fe
	google.golang.org/grpc v1.37.0
	google.golang.org/protobuf v1.26.0-rc.1
	mosn.io/layotto/sdk/go-sdk v0.0.0-20210604030858-e880c02bcf44
	mosn.io/layotto/spec v0.0.0-20210707123820-584778d048d3
)

replace (
	github.com/tetratelabs/proxy-wasm-go-sdk => github.com/layotto/proxy-wasm-go-sdk v0.14.1-0.20210929091432-0e4ff35b75af
	mosn.io/layotto/sdk/go-sdk => ../sdk/go-sdk
	mosn.io/layotto/spec => ../spec
)
