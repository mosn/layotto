// Copyright 2021 Layotto Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package utils

import (
	"errors"

	"github.com/hashicorp/consul/api"

	"mosn.io/layotto/components/lock"
)

type ConsulClient interface {
	Session() *api.Session
}

type ConsulKV interface {
	Acquire(p *api.KVPair, q *api.WriteOptions) (bool, *api.WriteMeta, error)
	Release(p *api.KVPair, q *api.WriteOptions) (bool, *api.WriteMeta, error)
}
type SessionFactory interface {
	Create(se *api.SessionEntry, q *api.WriteOptions) (string, *api.WriteMeta, error)
	Destroy(id string, q *api.WriteOptions) (*api.WriteMeta, error)
}

const (
	consulAddress  = "address"
	scheme         = "scheme"
	consulUsername = "username"
	consulPassword = "password"
	defaultScheme  = "http"
)

type ConsulMetadata struct {
	Address  string
	Scheme   string
	Username string
	Password string
}

func ParseConsulMetadata(meta lock.Metadata) (ConsulMetadata, error) {
	m := ConsulMetadata{}

	if val, ok := meta.Properties[consulAddress]; ok && val != "" {
		m.Address = val
	} else {
		return m, errors.New("consul error: missing host address")
	}

	m.Scheme = defaultScheme
	if val, ok := meta.Properties[scheme]; ok && val != "" {
		m.Scheme = val
	}

	if val, ok := meta.Properties[consulUsername]; ok && val != "" {
		m.Username = val
	}
	if val, ok := meta.Properties[consulPassword]; ok && val != "" {
		m.Password = val
	}

	return m, nil
}
