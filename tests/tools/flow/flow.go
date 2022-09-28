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

// Package flow is to help autoload env runners and run test cases.
package flow

import (
	"testing"

	"mosn.io/layotto/tests/tools/envs"
)

type testcase struct {
	name string
	fn   func(t *testing.T)
}

type Client struct {
	// runners is to setup the environment.
	runners envs.Runners

	// testcase related.
	testcases  []testcase
	setupFn    func(t *testing.T, name string)
	teardownFn func(t *testing.T, name string)
}

// New with runners. They will be executed Up with forwarding order and Down
// with reverse order.
//
//	client := New(runner1, runner2, runner3)
func New(runners ...envs.Runner) *Client {
	return &Client{
		runners:    runners,
		testcases:  make([]testcase, 0, 32),
		setupFn:    func(t *testing.T, name string) {},
		teardownFn: func(t *testing.T, name string) {},
	}
}

// Case of test.
func (c *Client) Case(name string, fn func(t *testing.T)) *Client {
	if fn != nil {
		c.testcases = append(c.testcases, testcase{name: name, fn: fn})
	}
	return c
}

// Setup for every case. You can determine which case is running through
// name.
func (c *Client) Setup(fn func(t *testing.T, name string)) *Client {
	if fn != nil {
		c.setupFn = fn
	}
	return c
}

// Teardown for every case. You can determine which case is running through
// name.
func (c *Client) Teardown(fn func(t *testing.T, name string)) *Client {
	if fn != nil {
		c.teardownFn = fn
	}
	return c
}

// Run all testcases.
func (c *Client) Run(t *testing.T) {
	if err := c.runE(t); err != nil {
		t.Errorf("[unittest] Run caught error: %+v", err)
	}
}

// runE is run all testcases that will return err.
func (c *Client) runE(t *testing.T) (err error) {
	if err = c.runners.Up(t); err != nil {
		return err
	}

	for _, testcase := range c.testcases {
		c.setupFn(t, testcase.name)
		t.Run(testcase.name, testcase.fn)
		c.teardownFn(t, testcase.name)
	}

	return c.runners.Down(t)
}
