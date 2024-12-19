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

package install

import (
	"context"
	"errors"

	wasm2 "mosn.io/mosn/pkg/wasm"

	"mosn.io/layotto/kit/logger"

	"mosn.io/layotto/pkg/filter/stream/common/http"
	"mosn.io/layotto/pkg/wasm"
)

func init() {
	w := wasm.GetDefault()
	w.AddEndpoint("install", NewEndpoint(w.Logger))
}

type Endpoint struct {
	logger logger.Logger
}

func NewEndpoint(log logger.Logger) *Endpoint {
	return &Endpoint{
		logger: log,
	}
}

func (e *Endpoint) Handle(ctx context.Context, params http.ParamsScanner) (map[string]interface{}, error) {
	conf, err := http.GetRequestData(ctx)
	if err != nil {
		e.logger.Errorf("[wasm][install] invalid request body for request /wasm/install, err:%v", err)
		return map[string]interface{}{"error": err.Error()}, err
	}

	if conf["name"] == nil {
		errorMessage := "can't get name property"
		e.logger.Errorf("[wasm][install] %v", errorMessage)
		return map[string]interface{}{"error": errorMessage}, errors.New(errorMessage)
	}

	id := conf["name"].(string)
	factory := wasm.GetFactory()
	if factory.IsRegister(id) {
		errorMessage := id + " is already registered"
		e.logger.Errorf("[wasm][install] %v", errorMessage)
		return map[string]interface{}{"error": errorMessage}, errors.New(errorMessage)
	}

	manager := wasm2.GetWasmManager()
	err = factory.Install(conf, manager)
	if err != nil {
		e.logger.Errorf("[wasm][install] %v", err)
		return map[string]interface{}{"error": err.Error()}, err
	}

	return nil, nil
}
