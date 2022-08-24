package lifecycle

import (
	"context"
	"sync"

	"mosn.io/layotto/components/pkg/common"
)

type ComponentKey struct {
	Kind string
	Name string
}

type dynamicComponentHolder struct {
	comp common.DynamicComponent
	mu   sync.Mutex
}

func (d *dynamicComponentHolder) ApplyConfig(ctx context.Context, metadata map[string]string) (err error) {
	// 1. lock
	d.mu.Lock()
	defer d.mu.Unlock()

	// 2. delegate to the comp
	return d.comp.ApplyConfig(ctx, metadata)
}

func ConcurrentDynamicComponent(comp common.DynamicComponent) common.DynamicComponent {
	return &dynamicComponentHolder{
		comp: comp,
		mu:   sync.Mutex{},
	}
}
