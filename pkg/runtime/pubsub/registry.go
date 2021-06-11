package pubsub

import (
	"fmt"
	"github.com/dapr/components-contrib/pubsub"
	"mosn.io/layotto/components/pkg/info"
)

const (
	ServiceName = "pubSub"
)

type Registry interface {
	Register(fs ...*Factory)
	Create(name string) (pubsub.PubSub, error)
}

type pubsubRegistry struct {
	stores map[string]func() pubsub.PubSub
	info   *info.RuntimeInfo
}

func NewRegistry(info *info.RuntimeInfo) Registry {
	info.AddService(ServiceName)
	return &pubsubRegistry{
		stores: make(map[string]func() pubsub.PubSub),
		info:   info,
	}
}

func (r *pubsubRegistry) Register(fs ...*Factory) {
	for _, f := range fs {
		r.stores[f.Name] = f.FactoryMethod
		r.info.RegisterComponent(ServiceName, f.Name)
	}
}

func (r *pubsubRegistry) Create(name string) (pubsub.PubSub, error) {
	if f, ok := r.stores[name]; ok {
		r.info.LoadComponent(ServiceName, name)
		return f(), nil
	}
	return nil, fmt.Errorf("service component %s is not regsitered", name)
}
