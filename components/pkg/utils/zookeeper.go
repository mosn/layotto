package utils

import (
	"errors"
	"fmt"
	"github.com/go-zookeeper/zk"
	"strconv"
	"strings"
	"time"
)

const (
	zkHost                = "zookeeperHosts"
	zkPassword            = "zookeeperPassword"
	sessionTimeout        = "SessionTimeout"
	logInfo               = "LogInfo"
	defaultSessionTimeout = 5 * time.Second
)

type ConnectionFactory interface {
	NewConnection(expire time.Duration, meta ZookeeperMetadata) (ZKConnection, error)
}

type ConnectionFactoryImpl struct {
}

func (c *ConnectionFactoryImpl) NewConnection(expire time.Duration, meta ZookeeperMetadata) (ZKConnection, error) {

	if expire == 0 {
		expire = meta.SessionTimeout
	}

	conn, _, err := zk.Connect(meta.Hosts, expire, zk.WithLogInfo(meta.LogInfo))
	if err != nil {
		return nil, err
	}
	return conn, nil
}

type ZKConnection interface {
	Get(path string) ([]byte, *zk.Stat, error)
	Set(path string, data []byte, version int32) (*zk.Stat, error)
	Delete(path string, version int32) error
	Create(path string, data []byte, flags int32, acl []zk.ACL) (string, error)
	Close()
}

type ZookeeperMetadata struct {
	Hosts          []string
	Password       string
	SessionTimeout time.Duration
	LogInfo        bool
}

func ParseZookeeperMetadata(properties map[string]string) (ZookeeperMetadata, error) {
	m := ZookeeperMetadata{}
	if val, ok := properties[zkHost]; ok && val != "" {
		split := strings.Split(val, ";")
		m.Hosts = append(m.Hosts, split...)
	} else {
		return m, errors.New("zookeeper store error: missing zkHost address")
	}

	if val, ok := properties[zkPassword]; ok && val != "" {
		m.Password = val
	}

	m.SessionTimeout = defaultSessionTimeout
	if val, ok := properties[sessionTimeout]; ok && val != "" {
		parsedVal, err := strconv.Atoi(val)
		if err != nil {
			return m, fmt.Errorf("zookeeper store error: can't parse SessionTimeout field: %s", err)
		}
		m.SessionTimeout = time.Duration(parsedVal) * time.Second
	}

	if val, ok := properties[logInfo]; ok && val != "" {
		b, err := strconv.ParseBool(val)
		if err != nil {
			return ZookeeperMetadata{}, err
		}
		m.LogInfo = b
	}
	return m, nil
}
