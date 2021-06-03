package pubsub

import "github.com/dapr/components-contrib/pubsub"

type Factory struct {
	Name          string
	FactoryMethod func() pubsub.PubSub
}

func NewFactory(name string, f func() pubsub.PubSub) *Factory {
	return &Factory{
		Name:          name,
		FactoryMethod: f,
	}
}
