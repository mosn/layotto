module github.com/layotto/layotto/demo

go 1.14

require (
	github.com/golang/protobuf v1.5.0
	github.com/layotto/layotto/sdk/go-sdk v0.0.0-20210604030858-e880c02bcf44
	github.com/layotto/layotto/spec v0.0.0-20210604023314-bb30491493a4
	google.golang.org/grpc v1.37.0
	google.golang.org/protobuf v1.26.0-rc.1
)

replace (
	github.com/layotto/layotto/sdk/go-sdk v0.0.0-20210604030858-e880c02bcf44 => ../sdk/go-sdk
	github.com/layotto/layotto/spec v0.0.0-20210604023314-bb30491493a4 => ../spec
)
