package actuators

import "sync"

type Indicator interface {
	Report() (status string, details map[string]interface{})
}

type ComponentsIndicator struct {
	ReadinessIndicator Indicator
	LivenessIndicator  Indicator
}

var componentsActutors sync.Map

func GetIndicatorWithName(name string) *ComponentsIndicator {
	if v, ok := componentsActutors.Load(name); ok {
		return v.(*ComponentsIndicator)
	}
	return nil
}

func SetComponentsActuators(name string, indicator *ComponentsIndicator) {
	componentsActutors.Store(name, indicator)
}
