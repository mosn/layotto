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
	"mosn.io/layotto/components/pkg/common"
	"testing"

	"mosn.io/layotto/components/hello"
)

func TestHelloWorld(t *testing.T) {
	hs := NewHelloWorld()
	hs.Init(&hello.HelloConfig{
		HelloString: "Hi",
	})

	req := &hello.HelloRequest{
		Name: "Layotto",
	}

	resp, _ := hs.Hello(context.Background(), req)
	if resp.HelloString != "Hi, Layotto" {
		t.Fatalf("hello output failed")
	}

	dc := hs.(common.DynamicComponent)
	dc.ApplyConfig(context.Background(), map[string]string{"hello": "Bye"})
	resp, _ = hs.Hello(context.Background(), req)
	if resp.HelloString != "Bye, Layotto" {
		t.Fatalf("hello output failed")
	}
}
