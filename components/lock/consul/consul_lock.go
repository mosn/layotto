package consul

import (
	"errors"
	"github.com/hashicorp/consul/api"
	"mosn.io/layotto/components/lock"
	"mosn.io/layotto/components/pkg/utils"
	"mosn.io/pkg/log"
	"strconv"
	"sync"
)

const (
	address       = "address"
	scheme        = "scheme"
	username      = "username"
	password      = "password"
	defaultScheme = "http"
)

type ConsulLock struct {
	metadata       metadata
	logger         log.ErrorLogger
	client         utils.ConsulClient
	sessionFactory utils.SessionFactory
	kv             utils.ConsulKV
	sMap           sync.Map
}

func NewConsulLock(logger log.ErrorLogger) *ConsulLock {
	consulLock := &ConsulLock{}
	return consulLock
}

func (c *ConsulLock) Init(metadata lock.Metadata) error {
	consulMetadata, err := parseConsulMetadata(metadata)
	if err != nil {
		return err
	}
	c.metadata = consulMetadata
	client, err := api.NewClient(&api.Config{
		Address: consulMetadata.address,
		Scheme:  consulMetadata.scheme,
	})
	c.client = client
	c.sessionFactory = client.Session()
	c.kv = client.KV()
	return nil
}
func (c *ConsulLock) Features() []lock.Feature {
	return nil
}

func getTTL(expire int32) string {
	//session TTL must be between [10s=24h0m0s]
	if expire < 10 {
		expire = 10
	}
	return strconv.Itoa(int(expire)) + "s"
}

func (c *ConsulLock) TryLock(req *lock.TryLockRequest) (*lock.TryLockResponse, error) {

	// create a session TTL
	session, _, err := c.sessionFactory.Create(&api.SessionEntry{
		TTL:       getTTL(req.Expire),
		LockDelay: 0,
		Behavior:  "delete", //Controls the behavior to delete when a session is invalidated.
	}, nil)

	if err != nil {
		return nil, err
	}

	// put a new KV pair with ttl session
	p := &api.KVPair{Key: req.ResourceId, Value: []byte(req.LockOwner), Session: session}
	//acquire lock
	acquire, _, err := c.kv.Acquire(p, nil)

	if err != nil {
		return nil, err
	}

	if acquire {
		//bind lockOwner+resourceId and session
		c.sMap.Store(req.LockOwner+"-"+req.ResourceId, session)
		return &lock.TryLockResponse{
			Success: true,
		}, nil
	} else {
		return &lock.TryLockResponse{
			Success: false,
		}, nil
	}

}
func (c *ConsulLock) Unlock(req *lock.UnlockRequest) (*lock.UnlockResponse, error) {

	session, ok := c.sMap.Load(req.LockOwner + "-" + req.ResourceId)

	if !ok {
		return &lock.UnlockResponse{Status: lock.LOCK_UNEXIST}, nil
	}
	// put a new KV pair with ttl session
	p := &api.KVPair{Key: req.ResourceId, Value: []byte(req.LockOwner), Session: session.(string)}
	//release lock
	release, _, err := c.kv.Release(p, nil)

	if err != nil {
		return &lock.UnlockResponse{Status: lock.INTERNAL_ERROR}, nil
	}

	if release {
		c.sMap.Delete(req.LockOwner + "-" + req.ResourceId)
		return &lock.UnlockResponse{Status: lock.SUCCESS}, nil
	} else {
		return &lock.UnlockResponse{Status: lock.LOCK_BELONG_TO_OTHERS}, nil
	}
}

type metadata struct {
	address  string
	scheme   string
	username string
	password string
}

func parseConsulMetadata(meta lock.Metadata) (metadata, error) {
	m := metadata{}

	if val, ok := meta.Properties[address]; ok && val != "" {
		m.address = val
	} else {
		return m, errors.New("consul error: missing host address")
	}

	m.scheme = defaultScheme
	if val, ok := meta.Properties[scheme]; ok && val != "" {
		m.scheme = val
	}

	if val, ok := meta.Properties[username]; ok && val != "" {
		m.username = val
	}
	if val, ok := meta.Properties[password]; ok && val != "" {
		m.password = val
	}

	return m, nil
}
