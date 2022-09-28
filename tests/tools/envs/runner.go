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
	"testing"
)

//go:generate mockgen -source=runner.go -destination=./mock/runner_mock.go -package=mock

// Runner is a executor to up & down the environment required for testing.
type Runner interface {
	// Up runner to setup environment.
	Up(tb testing.TB) error

	// Down runner to teardown envirment.
	Down(tb testing.TB) error
}

// Runners is alias of array of Runner and support functions to operation
// easily.
//
//	runners := New()
//	runners = runners.Add(runner1, runner2)
//	err := runners.Up(t)
//	assert.Nil(t, err)
//	err = runners.Down(t)
//	assert.Nil(t, err)
type Runners []Runner

func New() Runners {
	return make(Runners, 0, 32)
}

// Add non-nil runners.
func (rs Runners) Add(runners ...Runner) Runners {
	for _, runner := range runners {
		if runner != nil {
			rs = append(rs, runner)
		}
	}
	return rs
}

// Up with forwarding order.
func (rs Runners) Up(tb testing.TB) error {
	for _, runner := range rs {
		if err := runner.Up(tb); err != nil {
			return err
		}
	}
	return nil
}

// Down with reverse order.
func (rs Runners) Down(tb testing.TB) error {
	for i := len(rs) - 1; i >= 0; i-- {
		if err := rs[i].Down(tb); err != nil {
			return err
		}
	}
	return nil
}
