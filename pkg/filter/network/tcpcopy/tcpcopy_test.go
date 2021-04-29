package tcpcopy

import (
	"encoding/json"
	"gitlab.alipay-inc.com/ant-mesh/runtime/pkg/filter/network/tcpcopy/strategy"
	"mosn.io/api"
	"mosn.io/mosn/pkg/types"
	"mosn.io/pkg/buffer"
	"testing"
)

func TestCreateTcpcopyFactory(t *testing.T) {
	data := "{\"port\":\"12220\"}"
	var rpcCfg map[string]interface{}
	json.Unmarshal([]byte(data), &rpcCfg)

	type args struct {
		cfg map[string]interface{}
	}
	tests := []struct {
		name string
		args args
		want api.NetworkFilterChainFactory
	}{
		{
			name: "TestCreateTcpcopyFactory",
			args: struct{ cfg map[string]interface{} }{cfg: nil},
		},
		{
			name: "TestCreateTcpcopyFactory",
			args: struct{ cfg map[string]interface{} }{cfg: rpcCfg},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := CreateTcpcopyFactory(tt.args.cfg)
			if err != nil {
				t.Errorf("CreateTcpcopyFactory() error = %v", err)
				return
			}
		})
	}
}

func Test_tcpcopyFactory_OnData_switch_off(t *testing.T) {
	strategy.DumpSwitch = false

	type fields struct {
		tcpcopy *tcpcopy
	}
	type args struct {
		data types.IoBuffer
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantRes api.FilterStatus
	}{
		{
			name:    "",
			fields:  struct{ tcpcopy *tcpcopy }{tcpcopy: nil},
			args:    struct{ data types.IoBuffer }{data: nil},
			wantRes: api.Continue,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &tcpcopyFactory{
				tcpcopy: tt.fields.tcpcopy,
			}
			if gotRes := f.OnData(tt.args.data); gotRes != tt.wantRes {
				t.Errorf("OnData() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func Test_tcpcopyFactory_OnData_success(t *testing.T) {
	strategy.DumpSwitch = true
	strategy.DumpSampleFlag = 1

	type fields struct {
		tcpcopy *tcpcopy
	}
	tcpcopy_value := tcpcopy{port: "12220"}

	type args struct {
		data types.IoBuffer
	}

	buffer := buffer.NewIoBufferString("test")

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantRes api.FilterStatus
	}{
		{
			name:    "",
			fields:  struct{ tcpcopy *tcpcopy }{tcpcopy: &tcpcopy_value},
			args:    struct{ data types.IoBuffer }{data: buffer},
			wantRes: api.Continue,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &tcpcopyFactory{
				tcpcopy: tt.fields.tcpcopy,
			}
			if gotRes := f.OnData(tt.args.data); gotRes != tt.wantRes {
				t.Errorf("OnData() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func Test_tcpcopyFactory_OnNewConnection(t *testing.T) {
	type fields struct {
		tcpcopy *tcpcopy
	}
	tests := []struct {
		name   string
		fields fields
		want   api.FilterStatus
	}{
		{
			name:   "Test_tcpcopyFactory_OnNewConnection",
			fields: struct{ tcpcopy *tcpcopy }{tcpcopy: nil},
			want:   api.Continue,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &tcpcopyFactory{
				tcpcopy: tt.fields.tcpcopy,
			}
			if got := f.OnNewConnection(); got != tt.want {
				t.Errorf("OnNewConnection() = %v, want %v", got, tt.want)
			}
		})
	}
}
