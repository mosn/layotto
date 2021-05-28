package actuator

import "mosn.io/pkg/log"

type Actuator struct {
	endpointRegistry map[string]Endpoint
}

func New() *Actuator {
	return &Actuator{
		endpointRegistry: make(map[string]Endpoint),
	}
}

func (act *Actuator) GetEndpoint(name string) (endpoint Endpoint, ok bool) {
	e, ok := act.endpointRegistry[name]
	return e, ok
}

func (act *Actuator) AddEndpoint(name string, ep Endpoint) {
	_, ok := act.endpointRegistry[name]
	if ok {
		log.DefaultLogger.Warnf("Duplicate Endpoint name:  %v !", name)
	}
	act.endpointRegistry[name] = ep
}
