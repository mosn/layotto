package actuator

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSetAppInfoSingleton(t *testing.T) {
	SetAppInfoSingleton(nil)
	info := GetAppInfoSingleton()
	assert.True(t, info.Name == "")
	assert.True(t, info.Version == "")

	appInfo := NewAppInfo()
	appInfo.Name = "Test"
	appInfo.Version = "66666"
	SetAppInfoSingleton(appInfo)
	appInfo.Version = "7777"
	infoInterface, _ := GetAppContributor().GetInfo()
	info = infoInterface.(AppInfo)
	assert.True(t, info != *appInfo)
	assert.True(t, info.Name == "Test")
	assert.True(t, info.Version == "66666")
}
