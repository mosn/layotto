package zookeeper

import (
	"errors"
	"fmt"
	"github.com/go-zookeeper/zk"
	"mosn.io/layotto/components/lock"
	"mosn.io/pkg/log"
	"mosn.io/pkg/utils"
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

type ConnectionFactory interface {
	NewConnection(expire time.Duration, meta metadata) (ZKConnection, error)
}

type ConnectionFactoryImpl struct {
}

func (c *ConnectionFactoryImpl) NewConnection(expire time.Duration, meta metadata) (ZKConnection, error) {
	conn, _, err := zk.Connect(meta.hosts, expire*time.Second, zk.WithLogger(defaultLogger{}))
	if err != nil {
		return nil, err
	}
	return conn, nil
}

type ZKConnection interface {
	Get(path string) ([]byte, *zk.Stat, error)
	Delete(path string, version int32) error
	Create(path string, data []byte, flags int32, acl []zk.ACL) (string, error)
	Close()
}

type ZookeeperLock struct {
	//trylock reestablish connection  every time
	factory ConnectionFactory
	//unlock reuse this conneciton
	unlockConn ZKConnection
	metadata   metadata
	logger     log.ErrorLogger
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

func (p *ZookeeperLock) Init(metadata lock.Metadata) error {

	m, err := parseZookeeperMetadata(metadata)
	if err != nil {
		return err
	}

	p.metadata = m
	p.factory = &ConnectionFactoryImpl{}

	//init unlock connection
	zkConn, err := p.factory.NewConnection(p.metadata.sessionTimeout*time.Second, p.metadata)
	if err != nil {
		return err
	}
	p.unlockConn = zkConn
	return nil
}

func (p *ZookeeperLock) Features() []lock.Feature {
	return nil
}
func (p *ZookeeperLock) TryLock(req *lock.TryLockRequest) (*lock.TryLockResponse, error) {

	conn, err := p.factory.NewConnection(time.Duration(req.Expire)*time.Second, p.metadata)
	if err != nil {
		return &lock.TryLockResponse{}, err
	}
	//1.create zk ephemeral node
	_, err = conn.Create("/"+req.ResourceId, []byte(req.LockOwner), zk.FlagEphemeral, zk.WorldACL(zk.PermAll))

	//2.1 create node fail ,indicates lock fail
	if err != nil {
		defer conn.Close()
		//the node exists,lock fail
		if err == zk.ErrNodeExists {
			return &lock.TryLockResponse{
				Success: false,
			}, nil
		}
		//other err
		return &lock.TryLockResponse{}, err
	}

	//2.2 create node success, asyn  to make sure zkclient alive for need time
	utils.GoWithRecover(func() {
		//can also
		//time.Sleep(time.Second * time.Duration(req.Expire))
		timeAfterTrigger := time.After(time.Second * time.Duration(req.Expire))
		<-timeAfterTrigger
		// make sure close connecion
		conn.Close()
	}, nil)

	return &lock.TryLockResponse{
		Success: true,
	}, nil

}
func (p *ZookeeperLock) Unlock(req *lock.UnlockRequest) (*lock.UnlockResponse, error) {

	conn := p.unlockConn

	path := "/" + req.ResourceId
	owner, state, err := conn.Get(path)

	if err != nil {
		//node does not exist, indicates this lock has expired
		if err == zk.ErrNoNode {
			return &lock.UnlockResponse{Status: lock.LOCK_UNEXIST}, nil
		}
		//other err
		return &lock.UnlockResponse{}, err
	}
	//node exist ,but owner not this, indicates this lock has occupied or wrong unlock
	if string(owner) != req.LockOwner {
		return &lock.UnlockResponse{Status: lock.LOCK_BELONG_TO_OTHERS}, nil
	}
	err = conn.Delete(path, state.Version)
	//owner is this, but delete fail
	if err != nil {
		// delete no node , indicates this lock has expired
		if err == zk.ErrNoNode {
			return &lock.UnlockResponse{Status: lock.LOCK_UNEXIST}, nil
			// delete version error , indicates this lock has occupied by others
		} else if err == zk.ErrBadVersion {
			return &lock.UnlockResponse{Status: lock.LOCK_BELONG_TO_OTHERS}, nil
			//other error
		} else {
			return &lock.UnlockResponse{}, err
		}
	}
	//delete success, unlock success
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
