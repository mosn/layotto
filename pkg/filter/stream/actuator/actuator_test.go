package actuator

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

type MockEndpoint struct {
}

func (m *MockEndpoint) Handle(ctx context.Context, params ParamsScanner) (map[string]interface{}, error) {
	return nil, err
}

func TestActuator(t *testing.T) {
	act := GetDefault()
	endpoint, ok := act.GetEndpoint("health")
	assert.False(t, ok)
	assert.Nil(t, endpoint)

	act.AddEndpoint("", nil)
	endpoint, ok = act.GetEndpoint("")
	assert.True(t, ok)
	assert.Nil(t, endpoint)

	ep := &MockEndpoint{}
	act.AddEndpoint("health", ep)
	endpoint, ok = act.GetEndpoint("health")
	assert.True(t, ok)
	assert.Equal(t, endpoint, ep)
}
