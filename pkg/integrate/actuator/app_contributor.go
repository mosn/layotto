package actuator

import (
	"mosn.io/layotto/pkg/actuator/info"
	"sync/atomic"
	"time"
)

type AppContributor info.Contributor

func GetAppContributor() AppContributor {
	return info.ContributorAdapter(func() (interface{}, error) {
		return GetAppInfoSingleton(), nil
	})
}

var appInfoSingleton atomic.Value

type AppInfo struct {
	// The name of the program. Defaults to path.Base(os.Args[0])
	Name string `json:"name"`
	// Version of the program
	Version string `json:"version"`
	// Compilation date
	Compiled time.Time `json:"compiled,omitempty"`
}

func NewAppInfo() *AppInfo {
	return &AppInfo{}
}

func GetAppInfoSingleton() AppInfo {
	info, ok := appInfoSingleton.Load().(AppInfo)
	if ok {
		return info
	}
	return AppInfo{}
}

func SetAppInfoSingleton(a *AppInfo) {
	if a == nil {
		appInfoSingleton.Store(AppInfo{})
		return
	}
	appInfoSingleton.Store(*a)
}
