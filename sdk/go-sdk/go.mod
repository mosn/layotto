module mosn.io/layotto/sdk/go-sdk

go 1.14

require (
	mosn.io/layotto/spec v0.0.0-20210604023314-bb30491493a4
	github.com/pkg/errors v0.9.1
	github.com/stretchr/testify v1.7.0
	google.golang.org/grpc v1.37.0
	google.golang.org/protobuf v1.26.0-rc.1
)

replace mosn.io/layotto/spec v0.0.0-20210604023314-bb30491493a4 => ../../spec
