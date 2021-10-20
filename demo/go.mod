module mosn.io/layotto/demo

go 1.14

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/golang/protobuf v1.5.0
	github.com/google/uuid v1.2.0
	github.com/tetratelabs/proxy-wasm-go-sdk v0.14.1-0.20210922004205-46e3ac3a25fe
	golang.org/x/net v0.0.0-20190404232315-eb5bcb51f2a3 // indirect
	golang.org/x/sys v0.0.0-20200116001909-b77594299b42 // indirect
	golang.org/x/text v0.3.2 // indirect
	google.golang.org/grpc v1.37.0
	google.golang.org/protobuf v1.26.0-rc.1
	mosn.io/layotto/sdk/go-sdk 8288a9bdcb5eba5950df92be3ff0d39999820208
	mosn.io/layotto/spec 8288a9bdcb5eba5950df92be3ff0d39999820208
)

replace (
	github.com/tetratelabs/proxy-wasm-go-sdk => github.com/layotto/proxy-wasm-go-sdk v0.14.1-0.20210929091432-0e4ff35b75af
	mosn.io/layotto/sdk/go-sdk => ../sdk/go-sdk
	mosn.io/layotto/spec => ../spec
)
