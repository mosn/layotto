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
	"encoding/json"
	"errors"
)

type ContextKeyRequestData struct {
}

type RequestHandler interface {
	GetEndpoint(name string) (endpoint Endpoint, ok bool)
}

func GetRequestData(ctx context.Context) (map[string]interface{}, error) {
	requestData := ctx.Value(ContextKeyRequestData{})
	if requestData == nil {
		return nil, errors.New("invalid request body")
	}
	conf := make(map[string]interface{})
	err := json.Unmarshal(requestData.([]byte), &conf)
	if err != nil {
		return nil, err
	}
	return conf, nil
}
