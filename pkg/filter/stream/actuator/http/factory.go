package http

import (
	"context"
	"mosn.io/api"
	"mosn.io/mosn/pkg/log"
)

func init() {
	api.RegisterStream("actuator_filter", CreateActuatorFilterFactory)
}

type ServiceFactory struct{}

func (f *ServiceFactory) CreateFilterChain(context context.Context, callbacks api.StreamFilterChainFactoryCallbacks) {
	filter := &DispatchFilter{}
	callbacks.AddStreamReceiverFilter(filter, api.BeforeRoute)
}

func CreateActuatorFilterFactory(cfg map[string]interface{}) (api.StreamFilterChainFactory, error) {
	log.DefaultLogger.Infof("[actuator] create filter factory")
	return &ServiceFactory{}, nil
}
