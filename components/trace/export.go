package trace

import (
	"sync"
	"sync/atomic"
)

type Exporter interface {
	ExportSpan(s *Span)
}

type exportersMap map[string]Exporter

var activeExporters []string

var (
	exporterMu sync.RWMutex
	exporters  atomic.Value
)

func SetActiveExporters(exporter []string) {
	activeExporters = exporter
}

func GetExporter(name string) Exporter {
	exporterMu.RLock()
	defer exporterMu.RUnlock()
	if v, ok := exporters.Load().(exportersMap); ok {
		if ex, ok := v[name]; ok {
			return ex
		}
	}
	return nil
}

func RegisterExporter(name string, e Exporter) {
	exporterMu.Lock()
	newExporters := make(exportersMap)
	if old, ok := exporters.Load().(exportersMap); ok {
		for k, v := range old {
			newExporters[k] = v
		}
	}
	newExporters[name] = e
	exporters.Store(newExporters)
	exporterMu.Unlock()
}

func UnregisterExporter(name string) {
	exporterMu.Lock()
	newExporters := make(exportersMap)
	if old, ok := exporters.Load().(exportersMap); ok {
		for k, v := range old {
			newExporters[k] = v
		}
	}
	delete(newExporters, name)
	exporters.Store(newExporters)
	exporterMu.Unlock()
}
