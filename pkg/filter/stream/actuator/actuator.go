package actuator

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
	act.endpointRegistry[name] = ep
}
