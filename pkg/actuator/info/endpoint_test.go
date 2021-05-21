package info

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

type MockContributor struct {
}

func (m MockContributor) GetInfo() (info interface{}, err error) {
	return map[string]string{
		"k":  "v",
		"k1": "v1",
	}, nil
}

func TestEndpoint_Handle(t *testing.T) {
	ep := NewEndpoint()
	handle, err := ep.Handle(context.Background(), nil)
	assert.True(t, err == nil)
	assert.True(t, len(handle) == 0)

	AddInfoContributor("test", nil)
	handle, err = ep.Handle(context.Background(), nil)
	assert.True(t, err == nil)
	assert.True(t, len(handle) == 0)

	AddInfoContributor("test", MockContributor{})
	handle, err = ep.Handle(context.Background(), nil)
	assert.True(t, err == nil)
	assert.True(t, len(handle) == 1)
	assert.True(t, handle["test"].(map[string]string)["k"] == "v")
	assert.True(t, handle["test"].(map[string]string)["k1"] == "v1")
}
