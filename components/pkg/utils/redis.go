package utils

import (
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"strconv"
	"time"
)

const (
	db                     = "db"
	host                   = "redisHost"
	password               = "redisPassword"
	enableTLS              = "enableTLS"
	maxRetries             = "maxRetries"
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
