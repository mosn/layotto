package common

import (
	"testing"
)

func TestGetFileSize(t *testing.T) {
	type args struct {
		f string
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		{
			name: "TestGetFileSize",
			args: struct{ f string }{f: "/home/admin/logs/mosn/default.log"},
			want: -1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetFileSize(tt.args.f); got != tt.want {
				t.Errorf("GetFileSize() = %v, want %v", got, tt.want)
			}
		})
	}
}
