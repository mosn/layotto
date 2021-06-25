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

package hello

const ServiceName = "hello"

type HelloService interface {
	Init(*HelloConfig) error
	Hello(*HelloRequest) (*HelloReponse, error)
}

type HelloConfig struct {
	HelloString string `json:"hello"`
}

type HelloRequest struct {
	Name string `json:"name"`
}

type HelloReponse struct {
	HelloString string `json:"hello"`
}
