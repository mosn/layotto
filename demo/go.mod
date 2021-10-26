module mosn.io/layotto/demo

go 1.14

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/golang/protobuf v1.5.0
	github.com/google/uuid v1.2.0
	github.com/tetratelabs/proxy-wasm-go-sdk v0.14.1-0.20210922004205-46e3ac3a25fe
	google.golang.org/grpc v1.37.0
	google.golang.org/protobuf v1.26.0-rc.1
	mosn.io/layotto/sdk/go-sdk v0.0.0-20211020084508-6f5ee3cfeba0
	mosn.io/layotto/spec v0.0.0-20211020084508-6f5ee3cfeba0
)

replace (
	github.com/tetratelabs/proxy-wasm-go-sdk => github.com/layotto/proxy-wasm-go-sdk v0.14.1-0.20210929091432-0e4ff35b75af
	mosn.io/layotto/sdk/go-sdk => ../sdk/go-sdk
	mosn.io/layotto/spec => ../spec
)
