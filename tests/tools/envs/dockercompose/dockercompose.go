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

// Package dockercompose is a runner to up & down environment through
// Docker-Compose.
package dockercompose

import (
	"testing"
)

type Runner struct {
	config   *Config
	filename string
}

// New docker-compose runner with filename of docker-compose.yml and options.
//
//	runner := New("./docker-compose.yml")
func New(filename string, opts ...Option) *Runner {
	config := &Config{
		envVas: make(map[string]string),
	}
	for _, opt := range opts {
		opt(config)
	}

	return &Runner{
		config:   config,
		filename: filename,
	}
}

func (r *Runner) Up(tb testing.TB) error {
	args := []string{
		"-f", r.filename,
		"up", "-d",
		"--remove-orphans",
	}
	return runDockerComposeCmd(tb, args, r.config.workingDir, r.config.envVas)
}

func (r *Runner) Down(tb testing.TB) error {
	args := []string{
		"-f", r.filename,
		"down", "-v",
	}

	return runDockerComposeCmd(tb, args, r.config.workingDir, r.config.envVas)
}
