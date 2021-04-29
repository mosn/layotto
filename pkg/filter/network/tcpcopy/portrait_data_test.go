package tcpcopy

import (
	"github.com/layotto/layotto/pkg/filter/network/tcpcopy/strategy"
	_type "github.com/layotto/layotto/pkg/filter/network/tcpcopy/type"
	"mosn.io/pkg/log"
	"testing"
)

func Test_isHandle_switch_off(t *testing.T) {
	log.DefaultLogger.SetLogLevel(log.DEBUG)

	strategy.DumpSwitch = false

	want := false
	if got := isHandle(_type.RPC); got != want {
		t.Errorf("isHandle() = %v, want %v", got, want)
	}
}

func Test_isHandle_mutex(t *testing.T) {
	log.DefaultLogger.SetLogLevel(log.DEBUG)

	strategy.DumpSwitch = true
	strategy.DumpSampleFlag = 1
	getAndSwapDumpBusinessCache(_type.RPC, 1)

	want := false
	if got := isHandle(_type.RPC); got != want {
		t.Errorf("isHandle() = %v, want %v", got, want)
	}
}

func Test_isHandle_normal(t *testing.T) {
	log.DefaultLogger.SetLogLevel(log.DEBUG)

	strategy.DumpSwitch = true
	strategy.DumpSampleFlag = 1
	strategy.DumpCpuMaxRate = 200
	strategy.DumpMemMaxRate = 200

	want := true
	if got := isHandle(_type.CONFIGURATION); got != want {
		t.Errorf("isHandle() = %v, want %v", got, want)
	}
}

func Test_getAndSwapDumpBusinessCache(t *testing.T) {
	log.DefaultLogger.SetLogLevel(log.DEBUG)

	getAndSwapDumpBusinessCache(_type.RPC, 0)
	want := 0
	if got := getAndSwapDumpBusinessCache(_type.RPC, 1); got != want {
		t.Errorf("isHandle() = %v, want %v", got, want)
	}
}

func TestUploadPortraitData_switch_off(t *testing.T) {

	log.DefaultLogger.SetLogLevel(log.DEBUG)

	strategy.DumpSwitch = false

	want := false
	if got := UploadPortraitData(_type.RPC, nil, nil); got != want {
		t.Errorf("UploadPortraitData() = %v, want %v", got, want)
	}
}

func TestUploadPortraitData_normal(t *testing.T) {

	log.DefaultLogger.SetLogLevel(log.DEBUG)

	strategy.DumpSwitch = true
	strategy.DumpSampleFlag = 1
	strategy.DumpCpuMaxRate = 200
	strategy.DumpMemMaxRate = 200
	data := "{\"key\":\"value\"}"

	want := true
	if got := UploadPortraitData(_type.STATE, data, nil); got != want {
		t.Errorf("UploadPortraitData() = %v, want %v", got, want)
	}
}
