module github.com/layotto/layotto/go-sdk

go 1.14

require (
	google.golang.org/grpc v1.37.0
	github.com/stretchr/testify v1.7.0
	github.com/pkg/errors v0.9.1
	github.com/layotto/layotto/spec v0.0.0-20210604023314-bb30491493a4
)

replace github.com/layotto/layotto/spec v0.0.0-20210604023314-bb30491493a4 => ../../spec
