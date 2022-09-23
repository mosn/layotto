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
	"net"
	"net/url"
	"os"
	"strconv"
	"testing"
	"time"

	"mosn.io/pkg/log"

	"mosn.io/layotto/components/sequencer"

	"github.com/stretchr/testify/assert"
	"go.etcd.io/etcd/server/v3/embed"
)

const key = "resource_xxx"

// GetFreePort returns a free port from the OS.
func GetFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}

func TestEtcd_Init(t *testing.T) {
	var err error
	var etcdServer *embed.Etcd
	var etcdTestDir = "init.test.etcd"
	port, _ := GetFreePort()
	var etcdUrl = "localhost:" + strconv.Itoa(port)

	etcdServer, err = startEtcdServer(etcdTestDir, port)
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
	port, _ := GetFreePort()
	var etcdUrl = "localhost:" + strconv.Itoa(port)

	etcdServer, err = startEtcdServer(etcdTestDir, port)
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

	support, _, err := comp.GetSegment(nil)
	assert.False(t, support)
	assert.Nil(t, err)
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
