package zookeeper

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-zookeeper/zk"
	"mosn.io/layotto/components/sequencer"
	"mosn.io/pkg/log"
	"strconv"
	"strings"
	"time"
)

const (
	host                  = "zookeeperHosts"
	password              = "zookeeperPassword"
	sessionTimeout        = "sessionTimeout"
	logInfo               = "logInfo"
	defaultSessionTimeout = 5 * time.Second
)

type ConnectionFactory interface {
	NewConnection(meta metadata) (ZKConnection, error)
}

type ConnectionFactoryImpl struct {
}

func (c *ConnectionFactoryImpl) NewConnection(meta metadata) (ZKConnection, error) {
	conn, _, err := zk.Connect(meta.hosts, meta.sessionTimeout, zk.WithLogInfo(meta.logInfo))
	if err != nil {
		return nil, err
	}
	return conn, nil
}

type ZKConnection interface {
	Set(path string, data []byte, version int32) (*zk.Stat, error)
	Get(path string) ([]byte, *zk.Stat, error)
	Close()
}

type ZookeeperSequencer struct {
	client   ZKConnection
	metadata metadata
	logger   log.ErrorLogger
	factory  ConnectionFactory
	ctx      context.Context
	cancel   context.CancelFunc
}

// NewZookeeperSequencer returns a new zookeeper sequencer
func NewZookeeperSequencer(logger log.ErrorLogger) *ZookeeperSequencer {
	s := &ZookeeperSequencer{
		logger: logger,
	}

	return s
}

func (s *ZookeeperSequencer) Init(config sequencer.Configuration) error {
	m, err := parseRedisMetadata(config)
	if err != nil {
		return err
	}
	//init
	s.metadata = m
	s.factory = &ConnectionFactoryImpl{}
	connection, err := s.factory.NewConnection(m)
	if err != nil {
		return err
	}
	s.client = connection
	s.ctx, s.cancel = context.WithCancel(context.Background())

	//check biggerThan
	for k, needV := range s.metadata.biggerThan {
		if needV <= 0 {
			continue
		}
		_, stat, err := s.client.Get("/" + k)
		if err != nil {
			//key not exist
			if err == zk.ErrNoNode {
				return fmt.Errorf("zookeeper sequencer error: can not satisfy biggerThan guarantee.key: %s, current key does not exist", k)
			}
			//other error
			return err
		}
		realV := int64(stat.Version)

		if realV < needV {
			return fmt.Errorf("zookeeper sequencer error: can not satisfy biggerThan guarantee.key: %s,current id:%v", k, realV)
		}

	}
	return err

}

func (s *ZookeeperSequencer) GetNextId(req *sequencer.GetNextIdRequest) (*sequencer.GetNextIdResponse, error) {

	stat, err := s.client.Set("/"+req.Key, []byte(""), -1)

	if err != nil {
		return nil, err
	}
	return &sequencer.GetNextIdResponse{
		NextId: int64(stat.Version),
	}, nil
}

func (s *ZookeeperSequencer) GetSegment(req *sequencer.GetSegmentRequest) (support bool, result *sequencer.GetSegmentResponse, err error) {
	return false, nil, nil
}
func (s *ZookeeperSequencer) Close() error {
	s.cancel()
	s.client.Close()
	return nil
}
func parseRedisMetadata(config sequencer.Configuration) (metadata, error) {
	m := metadata{}

	m.biggerThan = config.BiggerThan

	if val, ok := config.Properties[host]; ok && val != "" {
		split := strings.Split(val, ";")
		m.hosts = append(m.hosts, split...)
	} else {
		return m, errors.New("zookeeper store error: missing host address")
	}

	if val, ok := config.Properties[password]; ok && val != "" {
		m.password = val
	}

	m.sessionTimeout = defaultSessionTimeout
	if val, ok := config.Properties[sessionTimeout]; ok && val != "" {
		parsedVal, err := strconv.Atoi(val)
		if err != nil {
			return m, fmt.Errorf("zookeeper store error: can't parse sessionTimeout field: %s", err)
		}
		m.sessionTimeout = time.Duration(parsedVal) * time.Second
	}

	if val, ok := config.Properties[logInfo]; ok && val != "" {
		b, err := strconv.ParseBool(val)
		if err != nil {
			return metadata{}, err
		}
		m.logInfo = b
	}
	return m, nil
}

type metadata struct {
	hosts          []string
	password       string
	sessionTimeout time.Duration
	logInfo        bool
	keyPrefix      string
	biggerThan     map[string]int64
}
