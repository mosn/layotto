package health

import (
	"context"
	"errors"
	"github.com/layotto/layotto/pkg/filter/stream/actuator"
)

const (
	health_key     = "health"
	liveness_key   = "liveness"
	readiness_key  = "readiness"
	status_key     = "status"
	components_key = "components"
)

var (
	invalidTypeError = errors.New("health type invalid.")
	serviceDownError = errors.New("service unavailable.")
)

func init() {
	actuator.GetDefault().AddEndpoint(health_key, NewEndpoint())
}

var type2Indicators = make(map[string]map[string]Indicator)

type Endpoint struct {
}

func NewEndpoint() *Endpoint {
	return &Endpoint{}
}

// Handle will check health status.The structure of the returned map is like:
//
// {
//
//  "status": "DOWN",
//
//  "components": {
//
//    "readinessProbe": {
//
//      "status": "DOWN"
//
//    }
//
//  }
//
// }
func (e *Endpoint) Handle(ctx context.Context, params actuator.ParamsScanner) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	// 1. validate params
	if params == nil || !params.HasNext() {
		return result, invalidTypeError
	}
	healthType := params.Next()
	m, ok := type2Indicators[healthType]
	if !ok || len(m) == 0 {
		return result, invalidTypeError
	}
	// 2. traverse the indicator chain
	result[status_key] = UP
	var resultErr error = nil
	components := make(map[string]Health)
	result[components_key] = components
	for k, idc := range m {
		data := idc.Report()
		components[k] = data
		if data.Status == DOWN {
			result[status_key] = DOWN
			resultErr = serviceDownError
		}
	}
	return result, resultErr
}

// AddLivenessIndicator register health.Indicator for liveness check.It's not concurrent-safe,so please invoke it ONLY in init method.
func AddLivenessIndicator(name string, idc Indicator) {
	addIndicator(liveness_key, name, idc)
}

// AddLivenessIndicatorFunc register health.Indicator for liveness check.Indicator.It's not concurrent-safe,so please invoke it ONLY in init method.
func AddLivenessIndicatorFunc(name string, f func() Health) {
	addIndicator(liveness_key, name, IndicatorAdapter(f))
}

// AddReadinessIndicator register health.Indicator for readiness check.Indicator.It's not concurrent-safe,so please invoke it ONLY in init method.
func AddReadinessIndicator(name string, idc Indicator) {
	addIndicator(readiness_key, name, idc)
}

// AddReadinessIndicatorFunc register health.Indicator for readiness check.Indicator.It's not concurrent-safe,so please invoke it ONLY in init method.
func AddReadinessIndicatorFunc(name string, f func() Health) {
	addIndicator(readiness_key, name, IndicatorAdapter(f))
}

func addIndicator(indicatorType string, name string, idc Indicator) {
	if idc == nil {
		return
	}
	if type2Indicators[indicatorType] == nil {
		type2Indicators[indicatorType] = make(map[string]Indicator)
	}
	type2Indicators[indicatorType][name] = idc
}
