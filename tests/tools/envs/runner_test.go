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

package envs

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"mosn.io/layotto/tests/tools/envs/mock"
)

func TestRunners_Add(t *testing.T) {
	t.Run("Add nil runner", func(t *testing.T) {
		runners := New()
		runners = runners.Add(nil)
		assert.Equal(t, 0, len(runners))
	})

	t.Run("Add runner success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mock := mock.NewMockRunner(ctrl)

		runners := New()
		runners = runners.Add(mock)
		assert.Equal(t, 1, len(runners))
	})
}

func TestRunners_Up(t *testing.T) {
	t.Run("Up failed", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockErr := errors.New("runner error")

		mock := mock.NewMockRunner(ctrl)
		mock.EXPECT().Up(gomock.Any()).Return(mockErr)

		err := New().Add(mock).Up(t)
		assert.Equal(t, mockErr, err)
	})

	t.Run("Up success with forwarding order", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		orders := []int{}

		mocks := make([]Runner, 0, 3)
		for i := 0; i < 3; i++ {
			num := i
			mock := mock.NewMockRunner(ctrl)
			mock.EXPECT().
				Up(gomock.Any()).
				DoAndReturn(func(_ testing.TB) error {
					orders = append(orders, num)
					return nil
				})
			mocks = append(mocks, mock)
		}

		err := New().Add(mocks...).Up(t)
		assert.Nil(t, err)
		assert.Equal(t, []int{0, 1, 2}, orders)
	})
}

func TestRunners_Down(t *testing.T) {
	t.Run("Down failed", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockErr := errors.New("runner error")

		mock := mock.NewMockRunner(ctrl)
		mock.EXPECT().Down(gomock.Any()).Return(mockErr)

		err := New().Add(mock).Down(t)
		assert.Equal(t, mockErr, err)
	})

	t.Run("Down success with reverse order", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		orders := []int{}

		mocks := make([]Runner, 0, 3)
		for i := 0; i < 3; i++ {
			num := i
			mock := mock.NewMockRunner(ctrl)
			mock.EXPECT().
				Down(gomock.Any()).
				DoAndReturn(func(_ testing.TB) error {
					orders = append(orders, num)
					return nil
				})
			mocks = append(mocks, mock)
		}
		err := New().Add(mocks...).Down(t)
		assert.Nil(t, err)
		assert.Equal(t, []int{2, 1, 0}, orders)
	})
}
