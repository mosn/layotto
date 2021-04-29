package client

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSayHello(t *testing.T) {
	item := &SayHelloRequest{"helloworld"}
	resp, err := testClient.SayHello(context.Background(), item)
	assert.Nil(t, err)
	assert.Equal(t, resp.Hello, "world")
}
