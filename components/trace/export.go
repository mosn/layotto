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
package trace

import (
	"sync"
)

// Exporter  is used to export Span
type Exporter interface {
	ExportSpan(s *Span)
}

var activeExporters []string

var (
	exporters sync.Map
)

func SetActiveExporters(exporter []string) {
	activeExporters = exporter
}

func GetExporter(name string) Exporter {
	if v, ok := exporters.Load(name); ok {
		return v.(Exporter)
	}
	return nil
}

func RegisterExporter(name string, e Exporter) {
	exporters.Store(name, e)
}

func UnregisterExporter(name string) {
	exporters.Delete(name)
	for i, v := range activeExporters {
		if v == name {
			activeExporters = append(activeExporters[:i], activeExporters[i+1:]...)
			return
		}
	}
}
