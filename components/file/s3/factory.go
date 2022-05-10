package s3

import "sync"

var initFuncRegistry map[string]S3ClientInit
var mux sync.RWMutex

type S3ClientInit interface {
	Init(staticConf map[string]string, DynConf map[string]string)
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
