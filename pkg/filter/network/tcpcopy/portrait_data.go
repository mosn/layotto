package tcpcopy

import (
	"context"
	"encoding/json"
	"github.com/layotto/layotto/pkg/filter/network/tcpcopy/model"
	"github.com/layotto/layotto/pkg/filter/network/tcpcopy/persistence"
	"github.com/layotto/layotto/pkg/filter/network/tcpcopy/strategy"
	_type "github.com/layotto/layotto/pkg/filter/network/tcpcopy/type"
	"mosn.io/api"
	"mosn.io/mosn/pkg/log"
	"mosn.io/mosn/pkg/types"
	"strconv"
	"sync"
	"sync/atomic"
)

var lock sync.Mutex

func isHandle(businessType _type.BusinessType) bool {
	// Determine whether to continue sampling
	if !persistence.IsPersistence() {
		return false
	}

	// The same business type, in the same sampling period, only accept one data report
	value := getAndSwapDumpBusinessCache(businessType, 1)
	if value == 0 && atomic.LoadInt32(&strategy.DumpSampleFlag) != 0 {
		return true
	}

	if log.DefaultLogger.GetLogLevel() >= log.DEBUG {
		log.DefaultLogger.Debugf("%s the business %s has already uploaded portrait data in the same sample duration.", model.LogDumpKey, businessType)
	}
	return false
}

func getAndSwapDumpBusinessCache(businessType _type.BusinessType, new int) int {

	lock.Lock()
	defer lock.Unlock()

	value, ok := strategy.DumpBusinessCache.LoadOrStore(businessType, new)
	if !ok {
		return 0 // 默认为0
	}

	tmp := value.(int)
	if tmp != new {
		strategy.DumpBusinessCache.Store(businessType, new)
		return tmp
	}

	return tmp
}

// Upload portrait data
func UploadPortraitData(businessType _type.BusinessType, data interface{}, ctx context.Context) bool {

	defer func() {
		if err := recover(); err != nil {
			log.DefaultLogger.Alertf(model.AlertDumpKey, "Upload portrait data error. %s", err)
		}
	}()

	if !isHandle(businessType) {
		if log.DefaultLogger.GetLogLevel() >= log.DEBUG {
			log.DefaultLogger.Debugf("%s ignore uploaded portrait data, condition does not match.", model.LogDumpKey)
		}
		return false
	}

	if log.DefaultLogger.GetLogLevel() >= log.DEBUG {
		log.DefaultLogger.Debugf("%s the uploaded portrait data is accepted.", model.LogDumpKey)
	}

	// Persistent user reported data
	var dataBytes []byte
	var err error = nil
	tmp := make(map[string]string)
	if _, ok := data.(api.HeaderMap); ok {
		data.(api.HeaderMap).Range(func(key, value string) bool {
			tmp[key] = value
			return true
		})
		dataBytes, err = json.Marshal(tmp)
	} else {
		dataBytes, err = json.Marshal(data)
	}

	if err != nil {
		log.DefaultLogger.Errorf("%s the uploaded portrait data is not json object.", model.LogDumpKey)
		return false
	}
	port := ""
	if ctx != nil {
		listener_port := ctx.Value(types.ContextKeyListenerPort)
		if listener_port != nil {
			port = strconv.Itoa(listener_port.(int))
		}
	}

	config := model.NewDumpUploadDynamicConfig(strategy.DumpSampleUuid, businessType, port, nil, string(dataBytes))
	persistence.GetDumpWorkPoolInstance().Schedule(config)

	return true
}
