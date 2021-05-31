package grpc

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewGRPCManager(t *testing.T) {
	m := NewManager()
	assert.NotNil(t, m)
}

func TestGetGRPCConnection(t *testing.T) {
	m := NewManager()
	assert.NotNil(t, m)
	port := 55555
	sslEnabled := false
	conn, err := m.GetGRPCConnection(fmt.Sprintf("127.0.0.1:%v", port), "", "", true, true, sslEnabled)
	assert.True(t, err != nil)
	assert.True(t, conn == nil)
}
