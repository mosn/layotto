module mosn.io/layotto/demo

go 1.14

require (
	github.com/golang/protobuf v1.5.0
	github.com/google/uuid v1.2.0
	github.com/minio/minio-go/v7 v7.0.15
	github.com/tetratelabs/proxy-wasm-go-sdk v0.14.1-0.20210922004205-46e3ac3a25fe
	golang.org/x/net v0.0.0-20210405180319-a5a99cb37ef4 // indirect
	golang.org/x/sys v0.0.0-20210510120138-977fb7262007 // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	google.golang.org/grpc v1.37.0
	google.golang.org/protobuf v1.26.0-rc.1
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
	mosn.io/layotto/sdk/go-sdk v0.0.0-20211020084508-6f5ee3cfeba0
	mosn.io/layotto/spec v0.0.0-20211020084508-6f5ee3cfeba0
)

replace (
	github.com/tetratelabs/proxy-wasm-go-sdk => github.com/layotto/proxy-wasm-go-sdk v0.14.1-0.20210929091432-0e4ff35b75af
	mosn.io/layotto/sdk/go-sdk => ../sdk/go-sdk
	mosn.io/layotto/spec => ../spec
)
