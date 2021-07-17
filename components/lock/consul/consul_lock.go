package consul

import (
	"errors"
	"github.com/hashicorp/consul/api"
	"mosn.io/layotto/components/lock"
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

type ConsulTxn interface {
	Txn(txn api.TxnOps, q *api.QueryOptions) (bool, *api.TxnResponse, *api.QueryMeta, error)
}
type SessionFactory interface {
	Create(se *api.SessionEntry, q *api.WriteOptions) (string, *api.WriteMeta, error)
}

type ConsulLock struct {
	metadata       metadata
	logger         log.ErrorLogger
	client         *api.Client
	sessionFactory SessionFactory
	txn            ConsulTxn
	mMap           map[string]lockMessage
	mutex          sync.Mutex
}

type lockMessage struct {
	LockOwner string
	Session   string
	Index     uint64
}

func NewConsulLock(logger log.ErrorLogger) *ConsulLock {
	consulLock := &ConsulLock{
		mMap: make(map[string]lockMessage, 0),
	}
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
	c.txn = client.Txn()
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
	}, nil)

	if err != nil {
		return nil, err
	}

	//use transaction to lock
	_, response, _, err := c.txn.Txn(api.TxnOps{
		&api.TxnOp{
			KV: &api.KVTxnOp{
				Verb:    api.KVLock,
				Key:     req.ResourceId,
				Value:   []byte(req.LockOwner),
				Session: session,
			},
		},
	}, nil)

	//lock successs
	if len(response.Errors) == 0 {
		c.mutex.Lock()
		defer c.mutex.Unlock()
		//renew map
		c.mMap[req.ResourceId] = lockMessage{
			LockOwner: req.LockOwner,
			Session:   session,
			Index:     response.Results[0].KV.ModifyIndex,
		}
		return &lock.TryLockResponse{
			Success: true,
		}, nil
		//lock fail
	} else {
		return &lock.TryLockResponse{
			Success: false,
		}, nil
	}
}
func (c *ConsulLock) Unlock(req *lock.UnlockRequest) (*lock.UnlockResponse, error) {

	c.mutex.Lock()
	defer c.mutex.Unlock()

	lockMes, ok := c.mMap[req.ResourceId]
	//lockMes exits but owner not this
	if ok {
		if lockMes.LockOwner != req.LockOwner {
			return &lock.UnlockResponse{Status: lock.LOCK_BELONG_TO_OTHERS}, nil
		}
		//lockMes not exits
	} else {
		return &lock.UnlockResponse{Status: lock.LOCK_UNEXIST}, nil
	}

	_, response, _, err := c.txn.Txn(api.TxnOps{
		&api.TxnOp{
			KV: &api.KVTxnOp{
				Verb:    api.KVDeleteCAS,
				Key:     req.ResourceId,
				Index:   lockMes.Index,
				Session: lockMes.Session,
			},
		},
	}, nil)

	if err != nil {
		return &lock.UnlockResponse{Status: lock.INTERNAL_ERROR}, nil
	}

	if len(response.Errors) == 0 {
		delete(c.mMap, req.ResourceId)
		return &lock.UnlockResponse{Status: lock.SUCCESS}, nil
		// this step contains preblem,consul cant distinct lock isn't held, or is held by another session
	} else {
		return &lock.UnlockResponse{Status: lock.LOCK_UNEXIST}, nil
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
