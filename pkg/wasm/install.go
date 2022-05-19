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

package wasm

import (
	"context"
	"encoding/json"
	"errors"
	"mosn.io/pkg/log"
)

func init() {
	GetDefault().AddEndpoint("install", &InstallEndpoint{})
}

type InstallEndpoint struct {
}

func (e *InstallEndpoint) Handle(_ context.Context, f *Filter) (map[string]interface{}, error) {
	conf := make(map[string]interface{})
	err := json.Unmarshal(f.receiverFilterHandler.GetRequestData().Bytes(), &conf)
	if err != nil {
		log.DefaultLogger.Errorf("[proxywasm][install] invalid body for request /wasm/install, err:%v", err)
		return nil, err
	}

	id := conf["name"]
	if id == nil {
		log.DefaultLogger.Errorf("[proxywasm][install] can't get name property")
		return nil, errors.New("can't get name property")
	}

	plugin, _ := f.router.GetRandomPluginByID(id.(string))
	if plugin != nil {
		log.DefaultLogger.Errorf("[proxywasm][install] %v is already registered", id)
		return nil, errors.New(id.(string) + " is already registered")
	}

	err = f.factory.register(conf)
	if err != nil {
		return nil, err
	}
	return nil, nil
}