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

package http

import (
	"context"

	"mosn.io/api"
	"mosn.io/mosn/pkg/log"
)

func RegisterFilter(filterType string, handler RequestHandler) {
	api.RegisterStream(filterType+"_filter", func(config map[string]interface{}) (api.StreamFilterChainFactory, error) {
		log.DefaultLogger.Infof("[%v] create filter factory", filterType)
		return &ServiceFactory{
			filterType:     filterType,
			requestHandler: handler,
		}, nil
	})
}

type ServiceFactory struct {
	filterType     string
	requestHandler RequestHandler
}

func (f *ServiceFactory) CreateFilterChain(context context.Context, callbacks api.StreamFilterChainFactoryCallbacks) {
	filter := &DispatchFilter{}
	filter.filterType = f.filterType
	filter.requestHandler = f.requestHandler
	callbacks.AddStreamReceiverFilter(filter, api.BeforeRoute)
}
