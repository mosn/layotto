package sequencer

import (
	"fmt"
	"mosn.io/layotto/components/pkg/info"
	"mosn.io/layotto/components/sequencer"
)

const (
	ServiceName = "sequencer"
)

type Registry interface {
	Register(fs ...*Factory)
	Create(name string) (sequencer.Store, error)
}

type sequencerRegistry struct {
	stores map[string]func() sequencer.Store
	info   *info.RuntimeInfo
}

func NewRegistry(info *info.RuntimeInfo) Registry {
	info.AddService(ServiceName)
	return &sequencerRegistry{
		stores: make(map[string]func() sequencer.Store),
		info:   info,
	}
}

func (r *sequencerRegistry) Register(fs ...*Factory) {
	for _, f := range fs {
		r.stores[f.Name] = f.FactoryMethod
		r.info.RegisterComponent(ServiceName, f.Name)
	}
}

func (r *sequencerRegistry) Create(name string) (sequencer.Store, error) {
	if f, ok := r.stores[name]; ok {
		r.info.LoadComponent(ServiceName, name)
		return f(), nil
	}
	return nil, fmt.Errorf("service component %s is not regsitered", name)
}
