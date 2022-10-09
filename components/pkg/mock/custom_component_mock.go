// Copyright 2021 Layotto Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package mock

import (
	"context"

	"mosn.io/layotto/components/custom"
)

type CustomComponentMock struct {
	ctx       context.Context
	config    *custom.Config
	initTimes int
}

func NewCustomComponentMock() custom.Component {
	return &CustomComponentMock{}
}

func (c *CustomComponentMock) InitTimes() int {
	return c.initTimes
}

func (c *CustomComponentMock) Initialize(ctx context.Context, config custom.Config) error {
	c.ctx = ctx
	c.config = &config
	c.initTimes++
	return nil
}

func (c *CustomComponentMock) GetReceivedConfig() *custom.Config {
	return c.config
}
func (c *CustomComponentMock) GetReceivedCtx() context.Context {
	return c.ctx
}
