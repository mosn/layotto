package sequencer

import (
	"mosn.io/layotto/components/sequencer"
)

type Factory struct {
	Name          string
	FactoryMethod func() sequencer.Store
}

func NewFactory(name string, f func() sequencer.Store) *Factory {
	return &Factory{
		Name:          name,
		FactoryMethod: f,
	}
}
