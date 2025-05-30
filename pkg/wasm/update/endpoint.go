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

package update

import (
	"context"
	"errors"

	"mosn.io/layotto/pkg/filter/stream/common/http"
	"mosn.io/layotto/pkg/wasm"

	wasm2 "mosn.io/mosn/pkg/wasm"

	"mosn.io/layotto/kit/logger"
)

func init() {
	w := wasm.GetDefault()
	w.AddEndpoint("update", NewEndpoint(w.Logger))
}

type Endpoint struct {
	logger logger.Logger
}

func NewEndpoint(logger logger.Logger) *Endpoint {
	return &Endpoint{
		logger: logger,
	}
}

func (e *Endpoint) Handle(ctx context.Context, params http.ParamsScanner) (map[string]interface{}, error) {
	conf, err := http.GetRequestData(ctx)
	if err != nil {
		e.logger.Errorf("[wasm][update] invalid request body for request /wasm/update, err:%v", err)
		return map[string]interface{}{"error": err.Error()}, err
	}

	if conf["name"] == nil {
		errorMessage := "can't get name property"
		e.logger.Errorf("[wasm][update] %v", errorMessage)
		return map[string]interface{}{"error": errorMessage}, errors.New(errorMessage)
	}

	if conf["instance_num"] == nil {
		errorMessage := "can't get instance_num property"
		e.logger.Errorf("[wasm][update] %v", errorMessage)
		return map[string]interface{}{"error": errorMessage}, errors.New(errorMessage)
	}

	instanceNum := int(conf["instance_num"].(float64))
	if instanceNum <= 0 {
		errorMessage := "instance_num should be greater than 0"
		e.logger.Errorf("[wasm][update] %v", errorMessage)
		return map[string]interface{}{"error": errorMessage}, errors.New(errorMessage)
	}

	id := (conf["name"]).(string)
	factory := wasm.GetFactory()
	err = factory.UpdateInstanceNum(id, instanceNum, wasm2.GetWasmManager())
	if err != nil {
		e.logger.Errorf("[wasm][update] %v", err)
		return map[string]interface{}{"error": err.Error()}, err
	}
	e.logger.Infof("[wasm] [update] wasm instance number updated success, id: %v, num: %v", id, instanceNum)
	return nil, nil
}
