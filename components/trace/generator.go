/*
 * Copyright 2021 Layotto Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package trace

import (
	"context"
	"sync"

	"mosn.io/api"
)

var (
	generators sync.Map
)

// Generator  is used to get or generate traceId/spanId/context
type Generator interface {
	Init(ctx context.Context)
	GetTraceId(ctx context.Context) string
	GetSpanId(ctx context.Context) string
	GenerateNewContext(ctx context.Context, span api.Span) context.Context
	GetParentSpanId(ctx context.Context) string
}

func RegisterGenerator(name string, ge Generator) {
	generators.Store(name, ge)
}

func GetGenerator(name string) Generator {
	g, ok := generators.Load(name)
	if ok {
		return g.(Generator)
	}
	return nil
}
