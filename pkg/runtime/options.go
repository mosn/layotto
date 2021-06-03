package runtime

import (
	"github.com/layotto/L8-components/configstores"
	"github.com/layotto/L8-components/hello"
	"github.com/layotto/layotto/pkg/services/pubsub"
	rgrpc "github.com/layotto/layotto/pkg/grpc"
	"github.com/layotto/layotto/pkg/services/rpc"
	"google.golang.org/grpc"
	"mosn.io/pkg/log"
)

// services encapsulates the service to include in the runtime
type services struct {
	hellos       []*hello.HelloFactory
	configStores []*configstores.StoreFactory
	rpcs         []*rpc.Factory
	pubSubs      []*pubsub.Factory
}

type runtimeOptions struct {
	// services
	services services
	// other config options
	srvMaker rgrpc.NewServer
	errInt   ErrInterceptor
	options  []grpc.ServerOption
}

type Option func(o *runtimeOptions)

func WithNewServer(f rgrpc.NewServer) Option {
	return func(o *runtimeOptions) {
		o.srvMaker = f
	}
}

func WithGrpcOptions(options ...grpc.ServerOption) Option {
	return func(o *runtimeOptions) {
		o.options = append(o.options, options...)
	}
}

type ErrInterceptor func(err error, format string, args ...interface{})

func WithErrInterceptor(i ErrInterceptor) Option {
	return func(o *runtimeOptions) {
		if o.errInt != nil {
			log.DefaultLogger.Fatalf("the error interceptor was already setted")
		}
		o.errInt = i
	}
}

// services options

func WithHelloFactory(hellos ...*hello.HelloFactory) Option {
	return func(o *runtimeOptions) {
		o.services.hellos = append(o.services.hellos, hellos...)
	}
}

func WithConfigStoresFactory(configStores ...*configstores.StoreFactory) Option {
	return func(o *runtimeOptions) {
		o.services.configStores = append(o.services.configStores, configStores...)
	}
}

func WithRpcFactory(rpcs ...*rpc.Factory) Option {
	return func(o *runtimeOptions) {
		o.services.rpcs = append(o.services.rpcs, rpcs...)
	}
}

func WithPubSubFactory(factorys ...*pubsub.Factory) Option {
	return func(o *runtimeOptions) {
		o.services.pubSubs = append(o.services.pubSubs, factorys...)
	}
}
