package redis

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"mosn.io/layotto/components/sequencer"
	"mosn.io/pkg/log"
	"strconv"
	"time"
)

const (
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

type StandaloneRedisSequencer struct {
	client   *redis.Client
	metadata metadata
	logger   log.ErrorLogger

	ctx    context.Context
	cancel context.CancelFunc
}

// EtcdSequencer returns a new etcd sequencer
func NewStandaloneRedisSequencer(logger log.ErrorLogger) *StandaloneRedisSequencer {
	s := &StandaloneRedisSequencer{
		logger: logger,
	}
	return s
}

func (s *StandaloneRedisSequencer) Init(config sequencer.Configuration) error {
	m, err := parseRedisMetadata(config)
	if err != nil {
		return err
	}
	//init
	s.metadata = m
	s.client = s.newClient(m)
	s.ctx, s.cancel = context.WithCancel(context.Background())

	//check biggerThan
	for k, bt := range s.metadata.biggerThan {
		if bt <= 0 {
			continue
		}
		get := s.client.Get(s.ctx, k)
		err := get.Err()
		if err != nil {
			//kv not exist
			if err == redis.Nil {
				return fmt.Errorf("standalone redis sequencer error: can not satisfy biggerThan guarantee.key: %s,current id:%v", k, 0)
			}
			//other error
			return err
		}
		val := get.Val()
		cur, err := strconv.ParseInt(val, 10, 64)
		//not int value
		if err != nil {
			return err
		}
		if cur < bt {
			return fmt.Errorf("standalone redis sequencer error: can not satisfy biggerThan guarantee.key: %s,current id:%v", k, cur)
		}
	}
	return err
}

func (p *StandaloneRedisSequencer) newClient(m metadata) *redis.Client {
	opts := &redis.Options{
		Addr:            m.host,
		Password:        m.password,
		DB:              defaultDB,
		MaxRetries:      m.maxRetries,
		MaxRetryBackoff: m.maxRetryBackoff,
	}
	if m.enableTLS {
		opts.TLSConfig = &tls.Config{
			InsecureSkipVerify: m.enableTLS,
		}
	}
	return redis.NewClient(opts)
}

func (s *StandaloneRedisSequencer) GetNextId(req *sequencer.GetNextIdRequest) (*sequencer.GetNextIdResponse, error) {

	incr := s.client.Incr(s.ctx, req.Key)

	err := incr.Err()
	if err != nil {
		return nil, err
	}

	return &sequencer.GetNextIdResponse{
		NextId: incr.Val(),
	}, nil
}

func (s *StandaloneRedisSequencer) GetSegment(req *sequencer.GetSegmentRequest) (support bool, result *sequencer.GetSegmentResponse, err error) {
	return false, nil, nil
}
func (s *StandaloneRedisSequencer) Close() error {
	s.cancel()
	return s.client.Close()
}
func parseRedisMetadata(config sequencer.Configuration) (metadata, error) {
	m := metadata{}

	m.biggerThan = config.BiggerThan

	if val, ok := config.Properties[host]; ok && val != "" {
		m.host = val
	} else {
		return m, errors.New("redis store error: missing host address")
	}

	if val, ok := config.Properties[password]; ok && val != "" {
		m.password = val
	}

	m.enableTLS = defaultEnableTLS
	if val, ok := config.Properties[enableTLS]; ok && val != "" {
		tls, err := strconv.ParseBool(val)
		if err != nil {
			return m, fmt.Errorf("redis store error: can't parse enableTLS field: %s", err)
		}
		m.enableTLS = tls
	}

	m.maxRetries = defaultMaxRetries
	if val, ok := config.Properties[maxRetries]; ok && val != "" {
		parsedVal, err := strconv.ParseInt(val, defaultBase, defaultBitSize)
		if err != nil {
			return m, fmt.Errorf("redis store error: can't parse maxRetries field: %s", err)
		}
		m.maxRetries = int(parsedVal)
	}

	m.maxRetryBackoff = defaultMaxRetryBackoff
	if val, ok := config.Properties[maxRetryBackoff]; ok && val != "" {
		parsedVal, err := strconv.ParseInt(val, defaultBase, defaultBitSize)
		if err != nil {
			return m, fmt.Errorf("redis store error: can't parse maxRetryBackoff field: %s", err)
		}
		m.maxRetryBackoff = time.Duration(parsedVal)
	}

	return m, nil
}

type metadata struct {
	host            string
	password        string
	maxRetries      int
	maxRetryBackoff time.Duration
	enableTLS       bool

	biggerThan map[string]int64
}
