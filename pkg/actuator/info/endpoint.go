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

	"mosn.io/pkg/log"

	"mosn.io/layotto/pkg/actuator"
	"mosn.io/layotto/pkg/filter/stream/common/http"
)

// init info Endpoint.
func init() {
	actuator.GetDefault().AddEndpoint("info", NewEndpoint())
}

var infoContributors = make(map[string]Contributor)

type Endpoint struct {
}

func NewEndpoint() *Endpoint {
	return &Endpoint{}
}

func (e *Endpoint) Handle(ctx context.Context, params http.ParamsScanner) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	var resultErr error
	// handle the infoContributors
	for k, c := range infoContributors {
		cinfo, err := c.GetInfo()
		if err != nil {
			log.DefaultLogger.Errorf("[actuator][info] Error when GetInfo.Contributor:%v,error:%v", k, err)
			result[k] = err.Error()
			resultErr = err
		} else {
			result[k] = cinfo
		}
	}
	return result, resultErr
}

// AddInfoContributor register info.Contributor.It's not concurrent-safe,so please invoke it ONLY in init method.
func AddInfoContributor(name string, c Contributor) {
	if c == nil {
		return
	}
	infoContributors[name] = c
}

// AddInfoContributorFunc register info.Contributor.It's not concurrent-safe,so please invoke it ONLY in init method.
func AddInfoContributorFunc(name string, f func() (interface{}, error)) {
	if f == nil {
		return
	}
	AddInfoContributor(name, ContributorAdapter(f))
}
