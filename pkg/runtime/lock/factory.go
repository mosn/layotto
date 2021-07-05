package lock

import "mosn.io/layotto/components/lock"

type Factory struct {
	Name          string
	FactoryMethod func() lock.LockStore
}

func NewFactory(name string, f func() lock.LockStore) *Factory {
	return &Factory{
		Name:          name,
		FactoryMethod: f,
	}
}
