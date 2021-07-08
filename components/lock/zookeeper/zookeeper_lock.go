package zookeeper

import (
	"errors"
	"fmt"
	"github.com/go-zookeeper/zk"
	"mosn.io/layotto/components/lock"
	"mosn.io/pkg/log"
	"strconv"
	"strings"
	"time"
)

const (
	host                  = "zookeeperHosts"
	password              = "zookeeperPassword"
	sessionTimeout        = "sessionTimeout"
	defaultSessionTimeout = 5
)

type ZKConnection interface {
	Get(path string) ([]byte, *zk.Stat, error)
	Delete(path string, version int32) error
	Create(path string, data []byte, flags int32, acl []zk.ACL) (string, error)
	Close()
	NewConnection(expire time.Duration, meta metadata) (ZKConnection, error)
}

type ZkConnectionImpl struct {
	connection *zk.Conn
}

func (z *ZkConnectionImpl) Get(path string) ([]byte, *zk.Stat, error) {
	return z.connection.Get(path)
}

func (z *ZkConnectionImpl) Delete(path string, version int32) error {
	return z.connection.Delete(path, version)
}
func (z *ZkConnectionImpl) Create(path string, data []byte, flags int32, acl []zk.ACL) (string, error) {
	return z.connection.Create(path, data, flags, acl)

}
func (z *ZkConnectionImpl) Close() {
	z.connection.Close()
}
func (z *ZkConnectionImpl) NewConnection(expire time.Duration, meta metadata) (ZKConnection, error) {

	conn, _, err := zk.Connect(meta.hosts, expire*time.Second, zk.WithLogger(defaultLogger{}))
	if err != nil {
		return nil, err
	}
	return &ZkConnectionImpl{connection: conn}, nil
}

type ZookeeperLock struct {
	conn     ZKConnection
	metadata metadata
	logger   log.ErrorLogger
}

func NewZookeeperLock(logger log.ErrorLogger) *ZookeeperLock {
	lock := &ZookeeperLock{
		logger: logger,
	}
	return lock
}

type defaultLogger struct{}

//silent the default logger
func (defaultLogger) Printf(format string, a ...interface{}) {

}

func (p *ZookeeperLock) newConnection(expire time.Duration) (ZKConnection, error) {
	//make sure a lock and a connection
	connection, err := p.conn.NewConnection(expire, p.metadata)

	if err != nil {
		return nil, err
	}
	return connection, nil
}

func (p *ZookeeperLock) Init(metadata lock.Metadata) error {

	m, err := parseZookeeperMetadata(metadata)
	if err != nil {
		return err
	}
	//nil to this
	p.conn = &ZkConnectionImpl{
		connection: nil,
	}
	p.metadata = m

	return nil
}

func (p *ZookeeperLock) Features() []lock.Feature {
	return nil
}
func (p *ZookeeperLock) TryLock(req *lock.TryLockRequest) (*lock.TryLockResponse, error) {

	conn, err := p.newConnection(time.Duration(req.Expire))
	if err != nil {
		return &lock.TryLockResponse{}, err
	}
	//1.create zk ephemeral node
	_, err = conn.Create("/"+req.ResourceId, []byte(req.LockOwner), zk.FlagEphemeral, zk.WorldACL(zk.PermAll))

	//2.1 create node fail ,indicates lock fail
	if err != nil {
		//make sure close connetion in time
		conn.Close()
		return &lock.TryLockResponse{
			Success: false,
		}, nil
	}

	//2.2 create node success, asyn  to make sure zkclient alive for need time
	go func() {
		//can also
		//time.Sleep(time.Second * time.Duration(req.Expire))
		timeAfterTrigger := time.After(time.Second * time.Duration(req.Expire))
		<-timeAfterTrigger
		// make sure close connecion
		conn.Close()
	}()

	return &lock.TryLockResponse{
		Success: true,
	}, nil

}
func (p *ZookeeperLock) Unlock(req *lock.UnlockRequest) (*lock.UnlockResponse, error) {

	conn, err := p.newConnection(p.metadata.sessionTimeout)
	defer conn.Close()

	if err != nil {
		return &lock.UnlockResponse{Status: lock.INTERNAL_ERROR}, err
	}

	path := "/" + req.ResourceId
	owner, state, err := conn.Get(path)
	//1. node does not exist, indicates this lock has expired or wrong unlock
	if err != nil {
		return &lock.UnlockResponse{Status: lock.LOCK_UNEXIST}, nil
	}
	//2. node exist ,but owner not this, indicates this lock has expired or wrong unlock
	if string(owner) != req.LockOwner {
		return &lock.UnlockResponse{Status: lock.LOCK_BELONG_TO_OTHERS}, nil
	}
	err = conn.Delete(path, state.Version)
	//3.owner is this, but delete fail, indicates this lock has expired
	//this step contains problem, the lock maybe also UNEXIST
	if err != nil {
		return &lock.UnlockResponse{Status: lock.LOCK_BELONG_TO_OTHERS}, nil
	}
	//4.delete success, unlock succes
	return &lock.UnlockResponse{Status: lock.SUCCESS}, nil
}

type metadata struct {
	hosts          []string
	password       string
	sessionTimeout time.Duration
}

func parseZookeeperMetadata(meta lock.Metadata) (metadata, error) {
	m := metadata{}
	if val, ok := meta.Properties[host]; ok && val != "" {
		split := strings.Split(val, ";")
		m.hosts = append(m.hosts, split...)
	} else {
		return m, errors.New("zookeeper store error: missing host address")
	}

	if val, ok := meta.Properties[password]; ok && val != "" {
		m.password = val
	}

	m.sessionTimeout = defaultSessionTimeout
	if val, ok := meta.Properties[sessionTimeout]; ok && val != "" {
		parsedVal, err := strconv.Atoi(val)
		if err != nil {
			return m, fmt.Errorf("zookeeper store error: can't parse sessionTimeout field: %s", err)
		}
		m.sessionTimeout = time.Duration(parsedVal)
	}
	return m, nil
}
