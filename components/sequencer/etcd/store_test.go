//
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
package etcd

import (
	"fmt"
	"mosn.io/layotto/components/sequencer"
	"mosn.io/pkg/log"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.etcd.io/etcd/server/v3/embed"
)

const key = "resource_xxx"
const key2 = "resource_xxx2"

const key3 = "resource_xxx3"
const key4 = "resource_xxx4"

func TestEtcd_Init(t *testing.T) {
	var err error
	var etcdServer *embed.Etcd
	var etcdTestDir = "init.test.etcd"
	var etcdUrl = "localhost:2380"

	etcdServer, err = startEtcdServer(etcdTestDir, 2380)
	assert.NoError(t, err)
	defer func() {
		etcdServer.Server.Stop()
		os.RemoveAll(etcdTestDir)
	}()
	comp := NewEtcdSequencer(log.DefaultLogger)

	cfg := sequencer.Configuration{
		BiggerThan: nil,
		Properties: map[string]string{},
	}
	err = comp.Init(cfg)
	assert.Error(t, err)

	cfg.Properties["endpoints"] = ""
	err = comp.Init(cfg)
	assert.Error(t, err)

	cfg.Properties["endpoints"] = etcdUrl
	cfg.Properties["dialTimeout"] = "a"
	err = comp.Init(cfg)
	assert.Error(t, err)

	cfg.Properties["dialTimeout"] = "2"
	err = comp.Init(cfg)
	assert.NoError(t, err)
	err = comp.Close()
	assert.NoError(t, err)

	//ca
	cfg.Properties["tlsCa"] = "/tmp"
	cfg.Properties["tlsCert"] = "/tmp"
	cfg.Properties["tlsCertKey"] = "/tmp"
	err = comp.Init(cfg)
	assert.Error(t, err)
}

func TestEtcd_CreateConnTimeout(t *testing.T) {
	var err error

	comp := NewEtcdSequencer(log.DefaultLogger)

	cfg := sequencer.Configuration{
		BiggerThan: nil,
		Properties: map[string]string{},
	}

	cfg.Properties["endpoints"] = "localhost:18888"
	cfg.Properties["dialTimeout"] = "3"
	startTime := time.Now()
	err = comp.Init(cfg)
	endTIme := time.Now()
	d := endTIme.Sub(startTime)
	assert.Error(t, err)
	assert.Equal(t, true, d.Seconds() >= 3)
}

func TestEtcd_GetNextId(t *testing.T) {
	var err error
	var resp *sequencer.GetNextIdResponse
	var etcdServer *embed.Etcd
	var etcdTestDir = "trylock.test.etcd"
	var etcdUrl = "localhost:23780"

	etcdServer, err = startEtcdServer(etcdTestDir, 23780)
	assert.NoError(t, err)
	defer func() {
		etcdServer.Server.Stop()
		os.RemoveAll(etcdTestDir)
	}()

	comp := NewEtcdSequencer(log.DefaultLogger)

	cfg := sequencer.Configuration{
		BiggerThan: nil,
		Properties: map[string]string{},
	}

	cfg.Properties["endpoints"] = etcdUrl
	err = comp.Init(cfg)
	assert.NoError(t, err)

	resp, err = comp.GetNextId(&sequencer.GetNextIdRequest{
		Key: key,
	})
	assert.NoError(t, err)
	var expected int64 = 1
	assert.Equal(t, expected, resp.NextId)

	//repeat
	resp, err = comp.GetNextId(&sequencer.GetNextIdRequest{
		Key: key,
	})
	assert.NoError(t, err)
	expected = 2
	assert.Equal(t, expected, resp.NextId)

}

func startEtcdServer(dir string, port int) (*embed.Etcd, error) {
	lc, _ := url.Parse(fmt.Sprintf("http://localhost:%v", port))
	lp, _ := url.Parse(fmt.Sprintf("http://localhost:%v", port+1))

	cfg := embed.NewConfig()
	cfg.Dir = dir
	cfg.LogLevel = "error"
	cfg.LCUrls = []url.URL{*lc}
	cfg.LPUrls = []url.URL{*lp}
	e, err := embed.StartEtcd(cfg)
	if err != nil {
		return nil, err
	}
	<-e.Server.ReadyNotify()
	return e, nil
}

func Test_addPathSeparator(t *testing.T) {
	p := addPathSeparator("")
	assert.Equal(t, p, "/")
	p = addPathSeparator("l8")
	assert.Equal(t, p, "/l8/")
	p = addPathSeparator("/l8")
	assert.Equal(t, p, "/l8/")
	p = addPathSeparator("l8/")
	assert.Equal(t, p, "/l8/")
}
