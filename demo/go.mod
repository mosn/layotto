module mosn.io/layotto/demo

go 1.14

require (
	github.com/golang/protobuf v1.5.3
	github.com/google/uuid v1.3.0
	github.com/gopherjs/gopherjs v0.0.0-20200217142428-fce0ec30dd00 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/minio/minio-go/v7 v7.0.15
	github.com/smartystreets/assertions v1.1.0 // indirect
	github.com/smartystreets/goconvey v1.6.6 // indirect
	google.golang.org/grpc v1.56.3
	google.golang.org/protobuf v1.30.0
	mosn.io/layotto/sdk/go-sdk v0.0.0-20211020084508-6f5ee3cfeba0
	mosn.io/layotto/spec v0.0.0-20220413092851-55c58dbb1a23
)

replace (
	mosn.io/layotto/sdk/go-sdk => ../sdk/go-sdk
	mosn.io/layotto/spec => ../spec

)
