package trace

import (
	"sync"
)

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
