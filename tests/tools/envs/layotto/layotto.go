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

// Packer layotto is to implement env runners to up & down the environments
// that test cases need.
package layotto

import (
	"testing"
	"time"

	"mosn.io/layotto/sdk/go-sdk/client"
	"mosn.io/layotto/tests/tools/envs/dockercompose"
)

type Runner struct {
	*dockercompose.Runner

	config *Config
	cli    client.Client
}

// NewE is to generate the Layotto env runner and return the error. It extends
// from Docker-Compose env runner.
func NewE(filename string, opts ...Option) (*Runner, error) {
	config := &Config{}
	for _, opt := range opts {
		opt(config)
	}

	dpRunner := dockercompose.New(filename,
		dockercompose.WithWorkingDir(config.workingDir),
		dockercompose.WithEnvVars(config.envVas))

	cli, err := client.NewClient()
	if err != nil {
		return nil, err
	}

	return &Runner{
		Runner: dpRunner,
		config: config,
		cli:    cli,
	}, nil
}

// New is a wrapper of NewE to avoid returning the error. It is friendly for
// coding.
func New(tb testing.TB, filename string, opts ...Option) *Runner {
	runner, err := NewE(filename, opts...)
	if err != nil {
		tb.Error(err)
	}
	return runner
}

func (r *Runner) Up(tb testing.TB) error {
	err := r.Runner.Up(tb)
	if err != nil {
		return err
	}
	return r.ensureStared()
}

// ensureStared is to ensure Laytto is started successfully.
func (r *Runner) ensureStared() error {

	// FIXME this is a temporary solution. Layotto have a healthcheck API,
	// however, it doesn't require all components implement this API. So we
	// don’t know if component is successed.
	time.Sleep(time.Second * 3)

	return nil
}
