package state

import (
	"fmt"
	"github.com/dapr/components-contrib/state"
	"mosn.io/layotto/components/pkg/info"
)

const (
	ServiceName = "state"
)

type Registry interface {
	Register(fs ...*Factory)
	Create(name string) (state.Store, error)
}

type stateRegistry struct {
	stores map[string]func() state.Store
	info   *info.RuntimeInfo
}

func NewRegistry(info *info.RuntimeInfo) Registry {
	info.AddService(ServiceName)
	return &stateRegistry{
		stores: make(map[string]func() state.Store),
		info:   info,
	}
}

func (r *stateRegistry) Register(fs ...*Factory) {
	for _, f := range fs {
		r.stores[f.Name] = f.FactoryMethod
		r.info.RegisterComponent(ServiceName, f.Name)
	}
}

func (r *stateRegistry) Create(name string) (state.Store, error) {
	if f, ok := r.stores[name]; ok {
		r.info.LoadComponent(ServiceName, name)
		return f(), nil
	}
	return nil, fmt.Errorf("service component %s is not regsitered", name)
}
