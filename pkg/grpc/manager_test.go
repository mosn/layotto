package grpc

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/connectivity"
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
	assert.NoError(t, err)
	conn2, err2 := m.GetGRPCConnection(fmt.Sprintf("127.0.0.1:%v", port), "", "", true, true, sslEnabled)
	assert.NoError(t, err2)
	assert.Equal(t, connectivity.Shutdown, conn.GetState())
	conn2.Close()
}
