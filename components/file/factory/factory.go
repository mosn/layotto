package factory

import (
	"encoding/json"
	"sync"
)

var initFuncRegistry map[string]S3ClientInit
var mux sync.RWMutex

type S3ClientInit func(staticConf json.RawMessage, dynConf map[string]string) (map[string]interface{}, error)

func init() {
	initFuncRegistry = make(map[string]S3ClientInit)
}
func RegisterInitFunc(name string, f S3ClientInit) {
	mux.Lock()
	initFuncRegistry[name] = f
	mux.Unlock()
}

func GetInitFunc(name string) S3ClientInit {
	mux.RLock()
	defer mux.RUnlock()
	return initFuncRegistry[name]
}
