package redis

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"mosn.io/layotto/components/lock"
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

// Standalone Redis lock store.Any fail-over related features are not supported,such as Sentinel and Redis Cluster.
type StandaloneRedisLock struct {
	client   *redis.Client
	metadata metadata
	replicas int

	features []lock.Feature
	logger   log.ErrorLogger

	ctx    context.Context
	cancel context.CancelFunc
}

// NewStandaloneRedisLock returns a new redis lock store
func NewStandaloneRedisLock(logger log.ErrorLogger) *StandaloneRedisLock {
	s := &StandaloneRedisLock{
		features: make([]lock.Feature, 0),
		logger:   logger,
	}

	return s
}

func (p *StandaloneRedisLock) Init(metadata lock.Metadata) error {
	// 1. parse config
	m, err := parseRedisMetadata(metadata)
	if err != nil {
		return err
	}
	p.metadata = m
	// 2. construct client
	p.client = p.newClient(m)
	p.ctx, p.cancel = context.WithCancel(context.Background())
	// 3. connect to redis
	if _, err = p.client.Ping(p.ctx).Result(); err != nil {
		return fmt.Errorf("[standaloneRedisLock]: error connecting to redis at %s: %s", m.host, err)
	}
	return err
}

func (p *StandaloneRedisLock) newClient(m metadata) *redis.Client {
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

func (p *StandaloneRedisLock) Features() []lock.Feature {
	return p.features
}

func (p *StandaloneRedisLock) TryLock(req *lock.TryLockRequest) (*lock.TryLockResponse, error) {
	nx := p.client.SetNX(p.ctx, req.ResourceId, req.LockOwner, time.Second*time.Duration(req.Expire))
	if nx == nil {
		return &lock.TryLockResponse{}, fmt.Errorf("[standaloneRedisLock]: SetNX returned nil.ResourceId: %s", req.ResourceId)
	}
	err := nx.Err()
	if err != nil {
		return &lock.TryLockResponse{}, err
	}

	return &lock.TryLockResponse{
		Success: nx.Val(),
	}, nil
}

const unlockScript = "local v = redis.call(\"get\",KEYS[1]); if v==false then return -1 end; if v~=ARGV[1] then return -2 else return redis.call(\"del\",KEYS[1]) end"

func (p *StandaloneRedisLock) Unlock(req *lock.UnlockRequest) (*lock.UnlockResponse, error) {
	// 1. delegate to client.eval lua script
	eval := p.client.Eval(p.ctx, unlockScript, []string{req.ResourceId}, req.LockOwner)
	// 2. check error
	if eval == nil {
		return newInternalErrorUnlockResponse(), fmt.Errorf("[standaloneRedisLock]: Eval unlock script returned nil.ResourceId: %s", req.ResourceId)
	}
	err := eval.Err()
	if err != nil {
		return newInternalErrorUnlockResponse(), err
	}
	// 3. parse result
	v := eval.Val()
	i, ok := v.(int)
	status := lock.INTERNAL_ERROR
	if ok {
		if i >= 0 {
			status = lock.SUCCESS
		} else if i == -1 {
			status = lock.LOCK_UNEXIST
		} else if i == -2 {
			status = lock.LOCK_BELONG_TO_OTHERS
		}
	}
	return &lock.UnlockResponse{
		Status: status,
	}, nil
}

func newInternalErrorUnlockResponse() *lock.UnlockResponse {
	return &lock.UnlockResponse{
		Status: lock.INTERNAL_ERROR,
	}
}

func (p *StandaloneRedisLock) Close() error {
	p.cancel()

	return p.client.Close()
}

func parseRedisMetadata(meta lock.Metadata) (metadata, error) {
	m := metadata{}

	if val, ok := meta.Properties[host]; ok && val != "" {
		m.host = val
	} else {
		return m, errors.New("redis store error: missing host address")
	}

	if val, ok := meta.Properties[password]; ok && val != "" {
		m.password = val
	}

	m.enableTLS = defaultEnableTLS
	if val, ok := meta.Properties[enableTLS]; ok && val != "" {
		tls, err := strconv.ParseBool(val)
		if err != nil {
			return m, fmt.Errorf("redis store error: can't parse enableTLS field: %s", err)
		}
		m.enableTLS = tls
	}

	m.maxRetries = defaultMaxRetries
	if val, ok := meta.Properties[maxRetries]; ok && val != "" {
		parsedVal, err := strconv.ParseInt(val, defaultBase, defaultBitSize)
		if err != nil {
			return m, fmt.Errorf("redis store error: can't parse maxRetries field: %s", err)
		}
		m.maxRetries = int(parsedVal)
	}

	m.maxRetryBackoff = defaultMaxRetryBackoff
	if val, ok := meta.Properties[maxRetryBackoff]; ok && val != "" {
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
}
