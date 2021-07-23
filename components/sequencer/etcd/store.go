package etcd

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"io/ioutil"
	"mosn.io/layotto/components/sequencer"
	"mosn.io/pkg/log"
	"strconv"
	"strings"
	"time"
)

const (
	defaultDialTimeout = 5
	defaultKeyPrefix   = "/layotto/"
	prefixKey          = "keyPrefixPath"
	usernameKey        = "username"
	passwordKey        = "password"
	dialTimeoutKey     = "dialTimeout"
	endpointsKey       = "endpoints"
	tlsCertPathKey     = "tlsCert"
	tlsCertKeyPathKey  = "tlsCertKey"
	tlsCaPathKey       = "tlsCa"
)

type EtcdSequencer struct {
	client   *clientv3.Client
	metadata metadata

	logger log.ErrorLogger

	ctx    context.Context
	cancel context.CancelFunc
}

// EtcdSequencer returns a new etcd sequencer
func NewEtcdSequencer(logger log.ErrorLogger) *EtcdSequencer {
	s := &EtcdSequencer{
		logger: logger,
	}

	return s
}

func (e *EtcdSequencer) Init(config sequencer.Configuration) error {
	// 1. parse config
	m, err := parseEtcdMetadata(config)
	if err != nil {
		return err
	}
	e.metadata = m
	// 2. construct client
	if e.client, err = e.newClient(m); err != nil {
		return err
	}
	e.ctx, e.cancel = context.WithCancel(context.Background())

	// 3. check biggerThan
	if len(e.metadata.biggerThan) > 0 {
		kv := clientv3.NewKV(e.client)
		for k, bt := range e.metadata.biggerThan {
			if bt <= 0 {
				continue
			}
			actualKey := e.getKeyInEtcd(k)
			get, err := kv.Get(e.ctx, actualKey)
			if err != nil {
				return err
			}
			var cur int64 = 0
			if get.Count > 0 && len(get.Kvs) > 0 {
				cur = get.Kvs[0].Version
			}
			if cur < bt {
				return fmt.Errorf("etcd sequencer error: can not satisfy biggerThan guarantee.key: %s,key in etcd: %s,current id:%v", k, actualKey, cur)
			}
		}
	}
	// TODO close component?
	return nil
}

func (e *EtcdSequencer) GetNextId(req *sequencer.GetNextIdRequest) (*sequencer.GetNextIdResponse, error) {
	key := e.getKeyInEtcd(req.Key)
	// Create new KV
	kv := clientv3.NewKV(e.client)
	// Create txn
	txn := kv.Txn(e.ctx)
	txn.If(clientv3.Compare(clientv3.CreateRevision(key), "=", 0)).Then(
		clientv3.OpPut(key, ""),
		clientv3.OpGet(key),
	).Else(
		clientv3.OpPut(key, ""),
		clientv3.OpGet(key),
	)
	// Commit
	txnResp, err := txn.Commit()
	if err != nil {
		return nil, err
	}
	return &sequencer.GetNextIdResponse{
		NextId: txnResp.Responses[1].GetResponseRange().Kvs[0].Version,
	}, nil
}

func (s *EtcdSequencer) GetSegment(req *sequencer.GetSegmentRequest) (support bool, result *sequencer.GetSegmentResponse, err error) {
	return false, nil, nil
}

func (e *EtcdSequencer) Close() error {
	e.cancel()

	return e.client.Close()
}

func (e *EtcdSequencer) newClient(meta metadata) (*clientv3.Client, error) {

	config := clientv3.Config{
		Endpoints:   meta.endpoints,
		DialTimeout: time.Second * time.Duration(meta.dialTimeout),
		Username:    meta.username,
		Password:    meta.password,
	}

	if meta.tlsCa != "" || meta.tlsCert != "" || meta.tlsCertKey != "" {
		//enable tls
		cert, err := tls.LoadX509KeyPair(meta.tlsCert, meta.tlsCertKey)
		if err != nil {
			return nil, fmt.Errorf("error reading tls certificate, cert: %s, certKey: %s, err: %s", meta.tlsCert, meta.tlsCertKey, err)
		}

		caData, err := ioutil.ReadFile(meta.tlsCa)
		if err != nil {
			return nil, fmt.Errorf("error reading tls ca %s, err: %s", meta.tlsCa, err)
		}

		pool := x509.NewCertPool()
		pool.AppendCertsFromPEM(caData)

		tlsConfig := &tls.Config{
			Certificates: []tls.Certificate{cert},
			RootCAs:      pool,
		}
		config.TLS = tlsConfig
	}

	if client, err := clientv3.New(config); err != nil {
		return nil, err
	} else {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(meta.dialTimeout))
		defer cancel()
		//ping
		_, err = client.Get(ctx, "ping")
		if err != nil {
			return nil, fmt.Errorf("etcd sequencer error: connect to etcd timeoout %s", meta.endpoints)
		}

		return client, nil
	}
}

func (e *EtcdSequencer) getKeyInEtcd(key string) string {
	return fmt.Sprintf("%s%s", e.metadata.keyPrefix, key)
}

func parseEtcdMetadata(config sequencer.Configuration) (metadata, error) {
	m := metadata{}
	var err error

	m.biggerThan = config.BiggerThan
	if val, ok := config.Properties[endpointsKey]; ok && val != "" {
		m.endpoints = strings.Split(val, ";")
	} else {
		return m, errors.New("etcd sequencer error: missing endpoints address")
	}

	if val, ok := config.Properties[dialTimeoutKey]; ok && val != "" {
		if m.dialTimeout, err = strconv.Atoi(val); err != nil {
			return m, fmt.Errorf("etcd sequencer error: ncorrect dialTimeout value %s", val)
		}
	} else {
		m.dialTimeout = defaultDialTimeout
	}

	if val, ok := config.Properties[prefixKey]; ok && val != "" {
		m.keyPrefix = addPathSeparator(val)
	} else {
		m.keyPrefix = defaultKeyPrefix
	}

	if val, ok := config.Properties[usernameKey]; ok && val != "" {
		m.username = val
	}

	if val, ok := config.Properties[passwordKey]; ok && val != "" {
		m.password = val
	}

	if val, ok := config.Properties[tlsCaPathKey]; ok && val != "" {
		m.tlsCa = val
	}

	if val, ok := config.Properties[tlsCertPathKey]; ok && val != "" {
		m.tlsCert = val
	}

	if val, ok := config.Properties[tlsCertKeyPathKey]; ok && val != "" {
		m.tlsCertKey = val
	}

	return m, nil
}

func addPathSeparator(p string) string {
	if p == "" {
		return "/"
	}
	if p[0] != '/' {
		p = "/" + p
	}
	if p[len(p)-1] != '/' {
		p = p + "/"
	}
	return p
}

type metadata struct {
	keyPrefix   string
	dialTimeout int
	endpoints   []string
	username    string
	password    string

	tlsCa      string
	tlsCert    string
	tlsCertKey string
	biggerThan map[string]int64
}
