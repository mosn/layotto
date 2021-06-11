package tcpcopy

import (
	"context"
	"encoding/json"
	"errors"
	"mosn.io/api"
	"mosn.io/layotto/pkg/filter/network/tcpcopy/model"
	"mosn.io/layotto/pkg/filter/network/tcpcopy/persistence"
	"mosn.io/layotto/pkg/filter/network/tcpcopy/strategy"
	v2 "mosn.io/mosn/pkg/config/v2"
	"mosn.io/mosn/pkg/types"
	"mosn.io/pkg/log"
	"net"
	"strconv"
)

func init() {
	api.RegisterNetwork("tcpcopy", CreateTcpcopyFactory)
}

var (
	ErrInvalidConfig = errors.New("invalid config for tcpcopy")
)

type config struct {
	port string
}

type tcpcopyFactory struct {
	cfg *config
}

func CreateTcpcopyFactory(cfg map[string]interface{}) (api.NetworkFilterChainFactory, error) {
	tcpConfig := &config{}
	// Parse static config for dump strategy
	if stg, ok := cfg["strategy"]; ok {
		data, err := json.Marshal(stg)
		if err != nil {
			log.DefaultLogger.Errorf("tcpcopy parse config error.%v", data)
		} else {
			strategy.UpdateAppDumpConfig(string(data))
		}
	}
	// TODO extract some other fields
	return &tcpcopyFactory{
		cfg: tcpConfig,
	}, nil
}

func (f *tcpcopyFactory) Init(param interface{}) error {
	// 1. get listener config
	cfg, ok := param.(*v2.Listener)
	if !ok {
		return ErrInvalidConfig
	}
	addr := cfg.AddrConfig
	if addr == "" {
		addr = cfg.Addr.String()
	}
	// 2. parse listener port
	var (
		netAddr *net.TCPAddr
		err     error
	)
	netAddr, err = net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		log.DefaultLogger.Errorf("invalid server address info: %s, error: %v", addr, err)
		return err
	}
	if netAddr.Port == 0 {
		log.DefaultLogger.Errorf("invalid server address info: %s", addr)
		return ErrInvalidConfig
	}
	// 3. set config
	f.cfg.port = strconv.Itoa(netAddr.Port)
	log.DefaultLogger.Debugf("tcpcopy filter initialized success")
	return nil
}

func (f *tcpcopyFactory) CreateFilterChain(context context.Context, callbacks api.NetWorkFilterChainFactoryCallbacks) {
	callbacks.AddReadFilter(f)
}

func (f *tcpcopyFactory) OnData(data types.IoBuffer) (res api.FilterStatus) {
	// Determine whether to continue sampling
	if !persistence.IsPersistence() {
		return api.Continue
	}

	// Asynchronous sampling
	config := model.NewDumpUploadDynamicConfig(strategy.DumpSampleUuid, "", f.cfg.port, data.Bytes(), "")
	persistence.GetDumpWorkPoolInstance().Schedule(config)
	return api.Continue
}

func (f *tcpcopyFactory) OnNewConnection() api.FilterStatus {
	return api.Continue
}

func (f *tcpcopyFactory) InitializeReadFilterCallbacks(cb api.ReadFilterCallbacks) {
}
