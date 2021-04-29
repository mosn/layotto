package common

import "testing"

func TestGetSystemUsageRate(t *testing.T) {
	_, _, err := GetSystemUsageRate()
	if err != nil {
		t.Errorf("GetSystemUsageRate() error = %v, wantErr %v", err, nil)
		return
	}
}
