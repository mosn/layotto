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
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

const (
	defaultKeyPrefix   = "/layotto/"
	defaultDialTimeout = 5
	prefixKey          = "keyPrefixPath"
	usernameKey        = "username"
	passwordKey        = "password"
	dialTimeoutKey     = "dialTimeout"
	endpointsKey       = "endpoints"
	tlsCertPathKey     = "tlsCert"
	tlsCertKeyPathKey  = "tlsCertKey"
	tlsCaPathKey       = "tlsCa"
)

func ParseEtcdMetadata(properties map[string]string) (EtcdMetadata, error) {
	m := EtcdMetadata{}
	var err error

	if val, ok := properties[endpointsKey]; ok && val != "" {
		m.Endpoints = strings.Split(val, ";")
	} else {
		return m, errors.New("etcd error: missing Endpoints address")
	}

	if val, ok := properties[dialTimeoutKey]; ok && val != "" {
		if m.DialTimeout, err = strconv.Atoi(val); err != nil {
			return m, fmt.Errorf("etcd error: ncorrect DialTimeout value %s", val)
		}
	} else {
		m.DialTimeout = defaultDialTimeout
	}

	if val, ok := properties[prefixKey]; ok && val != "" {
		m.KeyPrefix = addPathSeparator(val)
	} else {
		m.KeyPrefix = defaultKeyPrefix
	}

	if val, ok := properties[usernameKey]; ok && val != "" {
		m.Username = val
	}

	if val, ok := properties[passwordKey]; ok && val != "" {
		m.Password = val
	}

	if val, ok := properties[tlsCaPathKey]; ok && val != "" {
		m.TlsCa = val
	}

	if val, ok := properties[tlsCertPathKey]; ok && val != "" {
		m.TlsCert = val
	}

	if val, ok := properties[tlsCertKeyPathKey]; ok && val != "" {
		m.TlsCertKey = val
	}

	return m, nil
}

type EtcdMetadata struct {
	KeyPrefix   string
	DialTimeout int
	Endpoints   []string
	Username    string
	Password    string

	TlsCa      string
	TlsCert    string
	TlsCertKey string
}

func addPathSeparator(p string) string {
	if p == "" {
		return "/"
	}
	if p[0] != '/' {
		p = "/" + p
	}
	if p[len(p)-1] != '/' {
		p = p + "/"
	}
	return p
}

func NewEtcdClient(meta EtcdMetadata) (*clientv3.Client, error) {
	config := clientv3.Config{
		Endpoints:   meta.Endpoints,
		DialTimeout: time.Second * time.Duration(meta.DialTimeout),
		Username:    meta.Username,
		Password:    meta.Password,
	}

	if meta.TlsCa != "" || meta.TlsCert != "" || meta.TlsCertKey != "" {
		//enable tls
		cert, err := tls.LoadX509KeyPair(meta.TlsCert, meta.TlsCertKey)
		if err != nil {
			return nil, fmt.Errorf("error reading tls certificate, cert: %s, certKey: %s, err: %s", meta.TlsCert, meta.TlsCertKey, err)
		}

		caData, err := ioutil.ReadFile(meta.TlsCa)
		if err != nil {
			return nil, fmt.Errorf("error reading tls ca %s, err: %s", meta.TlsCa, err)
		}

		pool := x509.NewCertPool()
		pool.AppendCertsFromPEM(caData)

		tlsConfig := &tls.Config{
			Certificates: []tls.Certificate{cert},
			RootCAs:      pool,
		}
		config.TLS = tlsConfig
	}
	client, err := clientv3.New(config)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(meta.DialTimeout))
	defer cancel()
	//ping
	_, err = client.Get(ctx, "ping")
	if err != nil {
		return nil, fmt.Errorf("etcd error: connect to etcd timeoout %s", meta.Endpoints)
	}

	return client, nil
}
