package state

import (
	"github.com/dapr/components-contrib/state"
)

type Factory struct {
	Name          string
	FactoryMethod func() state.Store
}

func NewFactory(name string, f func() state.Store) *Factory {
	return &Factory{
		Name:          name,
		FactoryMethod: f,
	}
}
