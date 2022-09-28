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

package envrunner

import (
	"testing"
	"time"

	"github.com/helbing/testtools/envs/dockercompose"

	"mosn.io/layotto/sdk/go-sdk/client"
)

type Layotto struct {
	*dockercompose.Runner

	config *Config
	cli    client.Client
}

// NewLayottoE is to generate the Layotto env runner and return the error. It
// extends from Docker-Compose env runner.
func NewLayottoE(filename string, opts ...Option) (*Layotto, error) {
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

	return &Layotto{
		Runner: dpRunner,
		config: config,
		cli:    cli,
	}, nil
}

// NewLayotto is a wrapper of NewLayottoE to avoid returning the error. It is
// friendly for coding.
func NewLayotto(tb testing.TB, filename string, opts ...Option) *Layotto {
	runner, err := NewLayottoE(filename, opts...)
	if err != nil {
		tb.Error(err)
	}
	return runner
}

func (l *Layotto) Up(tb testing.TB) error {
	err := l.Runner.Up(tb)
	if err != nil {
		return err
	}
	return l.ensureStared()
}

// ensureStared is to ensure Laytto is started successfully.
func (l *Layotto) ensureStared() error {

	// FIXME this is a temporary solution. Layotto have a healthcheck API,
	// however, it doesn't require all components implement this API. So we
	// don’t know component is successed.
	time.Sleep(time.Second * 3)

	return nil
}
