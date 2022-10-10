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
	"crypto/tls"
	"errors"
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
)

const (
	db                     = "db"
	host                   = "redisHost"
	redisHosts             = "redisHosts"
	password               = "redisPassword"
	enableTLS              = "enableTLS"
	maxRetries             = "maxRetries"
	concurrency            = "concurrency"
	maxRetryBackoff        = "maxRetryBackoff"
	defaultBase            = 10
	defaultBitSize         = 0
	defaultDB              = 0
	defaultMaxRetries      = 3
	defaultMaxRetryBackoff = time.Second * 2
	defaultEnableTLS       = false
)

func NewRedisClient(m RedisMetadata) *redis.Client {
	opts := &redis.Options{
		Addr:            m.Host,
		Password:        m.Password,
		DB:              m.DB,
		MaxRetries:      m.MaxRetries,
		MaxRetryBackoff: m.MaxRetryBackoff,
	}
	if m.EnableTLS {
		opts.TLSConfig = &tls.Config{
			InsecureSkipVerify: m.EnableTLS,
		}
	}
	return redis.NewClient(opts)
}

type RedisMetadata struct {
	Host            string
	Password        string
	MaxRetries      int
	MaxRetryBackoff time.Duration
	EnableTLS       bool
	DB              int
}

func ParseRedisMetadata(properties map[string]string) (RedisMetadata, error) {
	m := RedisMetadata{}

	if val, ok := properties[host]; ok && val != "" {
		m.Host = val
	} else {
		return m, errors.New("redis store error: missing host address")
	}

	if val, ok := properties[password]; ok && val != "" {
		m.Password = val
	}

	m.EnableTLS = defaultEnableTLS
	if val, ok := properties[enableTLS]; ok && val != "" {
		tls, err := strconv.ParseBool(val)
		if err != nil {
			return m, fmt.Errorf("redis store error: can't parse enableTLS field: %s", err)
		}
		m.EnableTLS = tls
	}

	m.MaxRetries = defaultMaxRetries
	if val, ok := properties[maxRetries]; ok && val != "" {
		parsedVal, err := strconv.ParseInt(val, defaultBase, defaultBitSize)
		if err != nil {
			return m, fmt.Errorf("redis store error: can't parse maxRetries field: %s", err)
		}
		m.MaxRetries = int(parsedVal)
	}

	m.MaxRetryBackoff = defaultMaxRetryBackoff
	if val, ok := properties[maxRetryBackoff]; ok && val != "" {
		parsedVal, err := strconv.ParseInt(val, defaultBase, defaultBitSize)
		if err != nil {
			return m, fmt.Errorf("redis store error: can't parse maxRetryBackoff field: %s", err)
		}
		m.MaxRetryBackoff = time.Duration(parsedVal)
	}

	if val, ok := properties[db]; ok && val != "" {
		parsedVal, err := strconv.Atoi(val)
		if err != nil {
			return m, fmt.Errorf("redis store error: can't parse db field: %s", err)
		}
		m.DB = parsedVal
	} else {
		m.DB = defaultDB
	}
	return m, nil
}

func NewClusterRedisClient(m RedisClusterMetadata) []*redis.Client {
	clients := make([]*redis.Client, 0, len(m.Hosts))
	for _, Host := range m.Hosts {
		opts := &redis.Options{
			Addr:            Host,
			Password:        m.Password,
			DB:              m.DB,
			MaxRetries:      m.MaxRetries,
			MaxRetryBackoff: m.MaxRetryBackoff,
		}
		if m.EnableTLS {
			opts.TLSConfig = &tls.Config{
				InsecureSkipVerify: m.EnableTLS,
			}
		}
		clients = append(clients, redis.NewClient(opts))
	}
	return clients
}

type RedisClusterMetadata struct {
	Hosts           []string
	Concurrency     int
	Password        string
	MaxRetries      int
	MaxRetryBackoff time.Duration
	EnableTLS       bool
	DB              int
}

func ParseRedisClusterMetadata(properties map[string]string) (RedisClusterMetadata, error) {
	m := RedisClusterMetadata{}
	val, ok := properties[redisHosts]
	if !ok || val == "" {
		return m, errors.New("redis store error: missing redisHosts address")
	}
	hosts := strings.Split(val, ",")
	m.Hosts = hosts

	if val, ok := properties[password]; ok && val != "" {
		m.Password = val
	}

	m.EnableTLS = defaultEnableTLS
	if val, ok := properties[enableTLS]; ok && val != "" {
		tls, err := strconv.ParseBool(val)
		if err != nil {
			return m, fmt.Errorf("redis store error: can't parse enableTLS field: %s", err)
		}
		m.EnableTLS = tls
	}

	m.MaxRetries = defaultMaxRetries
	if val, ok := properties[maxRetries]; ok && val != "" {
		parsedVal, err := strconv.ParseInt(val, defaultBase, defaultBitSize)
		if err != nil {
			return m, fmt.Errorf("redis store error: can't parse maxRetries field: %s", err)
		}
		m.MaxRetries = int(parsedVal)
	}

	m.MaxRetryBackoff = defaultMaxRetryBackoff
	if val, ok := properties[maxRetryBackoff]; ok && val != "" {
		parsedVal, err := strconv.ParseInt(val, defaultBase, defaultBitSize)
		if err != nil {
			return m, fmt.Errorf("redis store error: can't parse maxRetryBackoff field: %s", err)
		}
		m.MaxRetryBackoff = time.Duration(parsedVal)
	}

	m.DB = defaultDB
	if val, ok := properties[db]; ok && val != "" {
		parsedVal, err := strconv.Atoi(val)
		if err != nil {
			return m, fmt.Errorf("redis store error: can't parse db field: %s", err)
		}
		m.DB = parsedVal
	}

	con, err := getConcurrency(properties)
	if err != nil {
		return m, err
	}
	m.Concurrency = con
	return m, nil
}

func getConcurrency(properties map[string]string) (int, error) {
	result := runtime.NumCPU()
	if val, ok := properties[concurrency]; ok && val != "" {
		con, err := strconv.Atoi(val)
		if err != nil {
			return result, fmt.Errorf("redis store error: can't parse concurrency field: %s", err)
		}
		if con > 0 {
			result = con
		}
	}
	return result, nil
}

func GetMiliTimestamp(i int64) int64 {
	return i / 1e6
}
