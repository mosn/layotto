// Copyright 2021 Layotto Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
