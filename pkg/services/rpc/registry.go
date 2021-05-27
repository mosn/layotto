package rpc

import (
	"fmt"
	"github.com/layotto/layotto/pkg/info"
)

const ServiceName = "rpc"

type Registry interface {
	Register(fs ...*Factory)
	Create(name string) (Invoker, error)
}

type RegistryImpl struct {
	// Key as implementing component name
	rpc  map[string]FactoryMethod
	info *info.RuntimeInfo
}

type FactoryMethod func() Invoker

func NewRegistry(info *info.RuntimeInfo) Registry {
	info.AddService(ServiceName)
	return &RegistryImpl{
		rpc:  make(map[string]FactoryMethod),
		info: info,
	}
}

func (r RegistryImpl) Register(fs ...*Factory) {
	for _, f := range fs {
		r.rpc[f.Name] = f.Fm
		r.info.RegisterComponent(ServiceName, f.Name)
	}
}

func (r RegistryImpl) Create(name string) (Invoker, error) {
	if f, ok := r.rpc[name]; ok {
		r.info.LoadComponent(ServiceName, name)
		return f(), nil
	}
	return nil, fmt.Errorf("service component %s is not registered", name)
}

type Factory struct {
	Name string
	Fm   FactoryMethod
}

func NewRpcFactory(name string, fm FactoryMethod) *Factory {
	return &Factory{
		Name: name,
		Fm:   fm,
	}
}
