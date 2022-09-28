/*
 * Copyright 2022 Layotto Authors
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

package flow

import (
	"errors"
	"reflect"
	"runtime"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"mosn.io/layotto/tests/tools/envs/mock"
)

func TestClient_Case(t *testing.T) {
	t.Run("Append case with nil", func(t *testing.T) {
		client := New()
		client.Case("mock name", nil).Case("mock name", nil)
		assert.Zero(t, len(client.testcases))
	})

	t.Run("Append case success", func(t *testing.T) {
		client := New()
		client.Case("mock name", func(t *testing.T) {})
		assert.Equal(t, 1, len(client.testcases))
	})
}

func TestClient_Setup(t *testing.T) {
	t.Run("Setup with nil", func(t *testing.T) {
		client := New()
		client.Setup(nil)
		assert.NotNil(t, client.setupFn)
	})

	t.Run("Setup success", func(t *testing.T) {
		client := New()
		fn := func(t *testing.T, name string) {}
		client.Setup(fn)
		assert.Same(t,
			runtime.FuncForPC(reflect.ValueOf(fn).Pointer()),
			runtime.FuncForPC(reflect.ValueOf(client.setupFn).Pointer()))
	})
}

func TestClient_Teardown(t *testing.T) {
	t.Run("Teardown with nil", func(t *testing.T) {
		client := New()
		client.Teardown(nil)
		assert.NotNil(t, client.teardownFn)
	})

	t.Run("Setup success", func(t *testing.T) {
		client := New()
		fn := func(t *testing.T, name string) {}
		client.Teardown(fn)
		assert.Same(t,
			runtime.FuncForPC(reflect.ValueOf(fn).Pointer()),
			runtime.FuncForPC(reflect.ValueOf(client.teardownFn).Pointer()))
	})
}

func TestClient_runE(t *testing.T) {
	t.Run("Without runners and tests", func(t *testing.T) {
		err := New().runE(t)
		assert.Nil(t, err)
	})

	t.Run("Runner up error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockErr := errors.New("runner error")

		mock := mock.NewMockRunner(ctrl)
		mock.EXPECT().Up(gomock.Any()).Return(mockErr)

		err := New(mock).runE(t)
		assert.Equal(t, mockErr, err)
	})

	t.Run("Runner down error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockErr := errors.New("runner error")

		mock := mock.NewMockRunner(ctrl)
		mock.EXPECT().Up(gomock.Any()).Return(nil)
		mock.EXPECT().Down(gomock.Any()).Return(mockErr)

		err := New(mock).runE(t)
		assert.Equal(t, mockErr, err)
	})

	t.Run("runE success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		callCount := 0
		setupCount := 0
		teardownCount := 0
		sum := 5

		mock := mock.NewMockRunner(ctrl)
		mock.EXPECT().Up(gomock.Any()).Return(nil)
		mock.EXPECT().Down(gomock.Any()).Return(nil)

		err := New(mock).
			Setup(func(t *testing.T, name string) {
				setupCount++
			}).
			Teardown(func(t *testing.T, name string) {
				teardownCount++
			}).
			Case("testcase1", func(t *testing.T) { callCount++ }).
			Case("testcase2", func(t *testing.T) { callCount++ }).
			Case("testcase3", func(t *testing.T) { callCount++ }).
			Case("testcase4", func(t *testing.T) { callCount++ }).
			Case("testcase5", func(t *testing.T) { callCount++ }).
			runE(t)

		assert.Nil(t, err)
		assert.Equal(t, sum, callCount)
		assert.Equal(t, sum, setupCount)
		assert.Equal(t, sum, teardownCount)
	})
}
