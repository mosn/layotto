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

package info

import (
	"context"
	"encoding/json"
	"errors"

	log "mosn.io/layotto/kit/logger"
	"mosn.io/layotto/pkg/actuator"
	"mosn.io/layotto/pkg/filter/stream/common/http"
)

// init info Endpoint.
func init() {
	actuator.GetDefault().AddEndpoint("logger", NewEndpoint())
}

type Endpoint struct {
}

type LoggerLevelChangedRequest struct {
	Component string `json:"component"`
	Level     string `json:"level"`
}

func NewEndpoint() *Endpoint {
	return &Endpoint{}
}

func (e *Endpoint) Handle(ctx context.Context, params http.ParamsScanner) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	requestData := ctx.Value(http.ContextKeyRequestData{})
	if requestData == nil {
		return nil, errors.New("invalid request body")
	}
	var request LoggerLevelChangedRequest
	err := json.Unmarshal(requestData.([]byte), &request)
	if err != nil {
		return nil, err
	}
	log.SetComponentLoggerLevel(request.Component, request.Level)
	var resultErr error
	// handle the infoContributors
	return result, resultErr
}
