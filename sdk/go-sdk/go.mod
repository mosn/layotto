module mosn.io/layotto/sdk/go-sdk

go 1.14

require (
	github.com/golang/protobuf v1.5.2
	github.com/google/uuid v1.3.0
	github.com/pkg/errors v0.9.1
	github.com/stretchr/testify v1.8.1
	google.golang.org/grpc v1.53.0
	google.golang.org/protobuf v1.28.1
	mosn.io/layotto/spec v0.0.0-20210707123820-584778d048d3
)

replace mosn.io/layotto/spec v0.0.0-20210707123820-584778d048d3 => ../../spec
