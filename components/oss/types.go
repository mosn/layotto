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

package oss

import (
	"encoding/json"
	"errors"

	"mosn.io/layotto/components/ref"
)

const (
	BasicConfiguration = "basic_config"
)

var (
	ErrInvalid = errors.New("invalid argument")
)

// Config wraps configuration for a oss implementation
type Config struct {
	ref.Config
	Metadata map[string]json.RawMessage `json:"metadata"`
	Type     string                     `json:"type"`
}
