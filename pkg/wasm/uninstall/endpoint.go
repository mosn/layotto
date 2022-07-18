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

package uninstall

import (
	"context"
	"errors"

	"mosn.io/layotto/pkg/wasm"

	"mosn.io/pkg/log"

	"mosn.io/layotto/pkg/filter/stream/common/http"
)

func init() {
	wasm.GetDefault().AddEndpoint("uninstall", NewEndpoint())
}

type Endpoint struct {
}

func NewEndpoint() *Endpoint {
	return &Endpoint{}
}

func (e *Endpoint) Handle(ctx context.Context, params http.ParamsScanner) (map[string]interface{}, error) {
	conf, err := http.GetRequestData(ctx)
	if err != nil {
		log.DefaultLogger.Errorf("[wasm][uninstall] invalid request body for request /wasm/uninstall, err:%v", err)
		return map[string]interface{}{"error": err.Error()}, err
	}

	if conf["name"] == nil {
		errorMessage := "can't get name property"
		log.DefaultLogger.Errorf("[wasm][uninstall] %v", errorMessage)
		return map[string]interface{}{"error": errorMessage}, errors.New(errorMessage)
	}

	factory := wasm.GetFactory()
	err = factory.UnInstall(conf["name"].(string))
	if err != nil {
		log.DefaultLogger.Errorf("[wasm][uninstall] %v", err)
		return map[string]interface{}{"error": err.Error()}, err
	}
	return nil, nil
}
