module mosn.io/layotto/demo

go 1.14

require (
	github.com/golang/protobuf v1.5.2
	github.com/google/uuid v1.3.0
	github.com/gopherjs/gopherjs v0.0.0-20200217142428-fce0ec30dd00 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/compress v1.13.6 // indirect
	github.com/kr/pretty v0.3.0 // indirect
	github.com/layotto/go-sdk v0.0.0-20241113124402-e55ee5816d2f
	github.com/minio/minio-go/v7 v7.0.15
	github.com/rogpeppe/go-internal v1.8.0 // indirect
	github.com/smartystreets/assertions v1.1.0 // indirect
	github.com/smartystreets/goconvey v1.6.6 // indirect
	golang.org/x/crypto v0.35.0 // indirect
	google.golang.org/genproto v0.0.0-20210602131652-f16073e35f0c // indirect
	google.golang.org/grpc v1.39.0
	google.golang.org/protobuf v1.27.1
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
	mosn.io/layotto/spec v0.0.0-20240927030843-b4fed4d06be4
)

replace mosn.io/layotto/spec => ../spec
