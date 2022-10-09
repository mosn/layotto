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
package lock

import (
	"strings"
	"testing"

	"mosn.io/layotto/components/lock"
	"mosn.io/layotto/components/pkg/info"
)

func TestNewRegistry(t *testing.T) {
	r := NewRegistry(info.NewRuntimeInfo())
	r.Register(NewFactory("mock", func() lock.LockStore {
		return nil
	}),
	)
	if _, err := r.Create("mock"); err != nil {
		t.Fatalf("create mock store failed: %v", err)
	}
	if _, err := r.Create("not exists"); !strings.Contains(err.Error(), "not regsitered") {
		t.Fatalf("create mock store failed: %v", err)
	}
}
