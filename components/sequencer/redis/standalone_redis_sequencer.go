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
	initNoOperation        = 0
	initReset              = 1
)

type StandaloneRedisSequencer struct {
	client   *redis.Client
	metadata metadata
	logger   log.ErrorLogger

	ctx    context.Context
	cancel context.CancelFunc
}

// NewStandaloneRedisSequencer returns a new redis sequencer
func NewStandaloneRedisSequencer(logger log.ErrorLogger) *StandaloneRedisSequencer {
	s := &StandaloneRedisSequencer{
		logger: logger,
	}
	return s
}

/*
   1. exists and >= biggerThan, no operation required, return 0
   2. not exists or < biggthan, reset val, return 1
   3. lua script occur error, such as tonumer(string), return error
*/
const initScript = `
if  redis.call('exists', KEYS[1])==1 and tonumber(redis.call('get', KEYS[1])) >= tonumber(ARGV[1]) then
    return 0
else
     redis.call('set', KEYS[1],ARGV[1])
     return 1
end
`

func (s *StandaloneRedisSequencer) Init(config sequencer.Configuration) error {
	m, err := parseRedisMetadata(config)
	if err != nil {
		return err
	}
	//init
	s.metadata = m
	s.client = s.newClient(m)
	s.ctx, s.cancel = context.WithCancel(context.Background())

	//check biggerThan, initialize if not satisfied
	for k, needV := range s.metadata.biggerThan {
		if needV <= 0 {
			continue
		}

		eval := s.client.Eval(s.ctx, initScript, []string{k}, needV)
		err = eval.Err()
		//occur error,  such as value is string type
		if err != nil {
			return err
		}
		//As long as there is no error, the initialization is successful
		//It may be a reset value or it may be satisfied before
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

func (s *StandaloneRedisSequencer) GetSegment(req *sequencer.GetSegmentRequest) (bool, *sequencer.GetSegmentResponse, error) {
	by := s.client.IncrBy(s.ctx, req.Key, int64(req.Size))
	err := by.Err()
	if err != nil {
		return true, nil, err
	}

	return true, &sequencer.GetSegmentResponse{
		From: by.Val() - int64(req.Size) + 1,
		To:   by.Val(),
	}, nil
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
