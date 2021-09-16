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

package client

// Copied from `mosn.io/pkg/utils/goroutine_test.go`.

import (
	"fmt"
	"io"
	"testing"
	"time"
)

func debugIgnoreLogger(w io.Writer, r interface{}) {
}

func TestGoWithRecover(t *testing.T) {
	panicStr := "test panic"
	panicHandler := func() {
		panic(panicStr)
	}
	output := ""
	recoverHandler := func(r interface{}) {
		output = fmt.Sprintf("%v", r)
	}
	GoWithRecover(panicHandler, recoverHandler)
	time.Sleep(time.Second) // wait panic goroutine
	if output != panicStr {
		t.Errorf("expected catch panic output, but got: %s", output)
	}
}

// recover handler panic, should not panic
func TestRecoverPanic(t *testing.T) {
	handler := func() {
		panic("1")
	}
	recoverHandler := func(r interface{}) {
		panic("2")
	}
	GoWithRecover(handler, recoverHandler)
}

// Example for how to recover with recover
type _run struct {
	count   int
	noPanic bool
}

func (r *_run) work() {
	GoWithRecover(r.exec, func(p interface{}) {
		r.work()
	})
}

func (r *_run) exec() {
	r.count++
	if r.count <= 2 {
		panic("panic")
	}
	r.noPanic = true
}

func TestGoWithRecoverAgain(t *testing.T) {
	r := &_run{}
	r.work()
	time.Sleep(time.Second)
	if !(r.noPanic && r.count == 3) {
		t.Errorf("panic handler is not restart expectedly, noPanic: %v, count: %d", r.noPanic, r.count)
	}
}
