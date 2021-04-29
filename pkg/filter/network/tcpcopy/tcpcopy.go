package tcpcopy

import (
	"context"
	"encoding/json"
	"github.com/layotto/layotto/pkg/filter/network/tcpcopy/model"
	"github.com/layotto/layotto/pkg/filter/network/tcpcopy/persistence"
	"github.com/layotto/layotto/pkg/filter/network/tcpcopy/strategy"
	"mosn.io/api"
	"mosn.io/mosn/pkg/types"
	"mosn.io/pkg/log"
)

func init() {
	api.RegisterNetwork("tcpcopy", CreateTcpcopyFactory)
}

type tcpcopy struct {
	port string
}

type tcpcopyFactory struct {
	tcpcopy *tcpcopy
}

func CreateTcpcopyFactory(cfg map[string]interface{}) (api.NetworkFilterChainFactory, error) {
	tcpconfig := &tcpcopy{}
	// Parse port number
	if portNum, ok := cfg["port"]; ok {
		tcpconfig.port = portNum.(string)
	}
	// Parse static config for dump strategy
	if stg, ok := cfg["strategy"]; ok {
		data, err := json.Marshal(stg)
		if err != nil {
			log.DefaultLogger.Errorf("tcpcopy parse config error.%v", data)
		} else {
			strategy.UpdateAppDumpConfig(string(data))
		}
	}

	return &tcpcopyFactory{
		tcpcopy: tcpconfig,
	}, nil
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
	config := model.NewDumpUploadDynamicConfig(strategy.DumpSampleUuid, "", f.tcpcopy.port, data.Bytes(), "")
	persistence.GetDumpWorkPoolInstance().Schedule(config)
	return api.Continue
}

func (f *tcpcopyFactory) OnNewConnection() api.FilterStatus {
	return api.Continue
}

func (f *tcpcopyFactory) InitializeReadFilterCallbacks(cb api.ReadFilterCallbacks) {
}
