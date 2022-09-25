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

package apollo

import (
	"fmt"

	"github.com/apolloconfig/agollo/v4"
	agolloConfig "github.com/apolloconfig/agollo/v4/env/config"
	"mosn.io/pkg/log"
)

// An interface to abstract different apollo sdks,also making it easier to write unit tests.
type Repository interface {
	SetConfig(r *repoConfig)
	Connect() error
	// subscribe
	AddChangeListener(listener *changeListener)
	// query
	Get(namespace string, key string) (interface{}, error)
	//	process every items under the namespace
	Range(namespace string, f func(key, value interface{}) bool) error
}

type repoConfig struct {
	addr          string
	appId         string
	storeName     string
	env           string
	cluster       string
	namespaceName string
	// whether backup config after fetch config from apollo
	isBackupConfig bool
	secret         string
}

func init() {
	agollo.SetLogger(NewDefaultLogger(log.DefaultLogger))
}

// Implement Repository interface
type AgolloRepository struct {
	client agollo.Client
	cfg    *repoConfig
}

func (a *AgolloRepository) Connect() error {
	var err error
	a.client, err = agollo.StartWithConfig(func() (*agolloConfig.AppConfig, error) {
		return repoConfig2AgolloConfig(a.cfg), nil
	})
	return err
}

func (a *AgolloRepository) SetConfig(r *repoConfig) {
	a.cfg = r
}

func repoConfig2AgolloConfig(r *repoConfig) *agolloConfig.AppConfig {
	return &agolloConfig.AppConfig{
		IP:             r.addr,
		AppID:          r.appId,
		Cluster:        r.cluster,
		NamespaceName:  r.namespaceName,
		IsBackupConfig: r.isBackupConfig,
		Secret:         r.secret,
	}
}

func newAgolloRepository() Repository {
	return &AgolloRepository{}
}

func (a *AgolloRepository) Get(namespace string, key string) (interface{}, error) {
	// 1. get cache
	cache := a.client.GetConfigCache(namespace)
	if cache == nil {
		return nil, fmt.Errorf("no cache for namespace:%v", namespace)
	}
	// 2. query value
	return cache.Get(key)
}

func (a *AgolloRepository) Range(namespace string, f func(key interface{}, value interface{}) bool) error {
	// 1. get cache
	cache := a.client.GetConfigCache(namespace)
	if cache == nil {
		return fmt.Errorf("no cache for namespace:%v", namespace)
	}
	// 2. loop process
	cache.Range(f)
	return nil
}

func (a *AgolloRepository) AddChangeListener(listener *changeListener) {
	a.client.AddChangeListener(listener)
}
