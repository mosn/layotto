module mosn.io/layotto/demo

go 1.14

require (
	github.com/golang/protobuf v1.5.0
	github.com/google/uuid v1.2.0
	github.com/minio/minio-go/v7 v7.0.15
	google.golang.org/grpc v1.37.0
	google.golang.org/protobuf v1.26.0-rc.1
	mosn.io/layotto/sdk/go-sdk v0.0.0-20211020084508-6f5ee3cfeba0
	mosn.io/layotto/spec v0.0.0-20211020084508-6f5ee3cfeba0
)

replace (
	mosn.io/layotto/sdk/go-sdk => ../sdk/go-sdk
	mosn.io/layotto/spec => ../spec
)
