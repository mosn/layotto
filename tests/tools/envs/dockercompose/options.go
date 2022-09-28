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

package dockercompose

type Config struct {
	workingDir string

	// envVas is environment variables.
	envVas map[string]string
}

// Option of New docker-compose runner.
type Option func(*Config)

func WithWorkingDir(dir string) Option {
	return func(c *Config) { c.workingDir = dir }
}

func WithEnvVars(envVas map[string]string) Option {
	return func(c *Config) { c.envVas = envVas }
}
