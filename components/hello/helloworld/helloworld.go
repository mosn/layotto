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

package helloworld

import (
	"context"
	"sync/atomic"

	"github.com/dapr/components-contrib/secretstores"

	"mosn.io/layotto/components/configstores"
	"mosn.io/layotto/components/hello"
)

type HelloWorld struct {
	Say         atomic.Value
	config      configstores.Store
	secretStore secretstores.SecretStore
}

func (hw *HelloWorld) ApplyConfig(ctx context.Context, metadata map[string]string) (err error) {
	greetings, ok := metadata["hello"]
	if !ok {
		return nil
	}
	hw.Say.Store(greetings)
	return nil
}

func (hw *HelloWorld) SetConfigStore(cs configstores.Store) (err error) {
	//save for use
	hw.config = cs
	return nil
}
func (hw *HelloWorld) SetSecretStore(ss secretstores.SecretStore) (err error) {
	//save for use
	hw.secretStore = ss
	return nil
}

var _ hello.HelloService = &HelloWorld{}

func NewHelloWorld() hello.HelloService {
	return &HelloWorld{}
}

func (hw *HelloWorld) Init(config *hello.HelloConfig) error {
	hw.Say.Store(config.HelloString)
	return nil
}

func (hw *HelloWorld) Hello(ctx context.Context, req *hello.HelloRequest) (*hello.HelloReponse, error) {
	greetings, _ := hw.Say.Load().(string)
	if req.Name != "" {
		greetings = greetings + ", " + req.Name
	}
	return &hello.HelloReponse{
		HelloString: greetings,
	}, nil
}
