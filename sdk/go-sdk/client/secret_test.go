package client

import (
	"context"
	"github.com/stretchr/testify/assert"
	runtimev1pb "mosn.io/layotto/spec/proto/runtime/v1"
	"testing"
)

var secretKey = "secret-key"
var secretValue = "secret-value"

func TestGetSecret(t *testing.T) {

	resp, err := testClient.GetSecret(context.Background(), &runtimev1pb.GetSecretRequest{
		Key: secretKey,
	})
	assert.Nil(t, err)
	assert.Equal(t, resp.Data[secretKey], secretValue)
}

func TestGetBulkSecret(t *testing.T) {
	resp, err := testClient.GetBulkSecret(context.Background(), &runtimev1pb.GetBulkSecretRequest{})
	assert.Nil(t, err)
	assert.Equal(t, resp.Data[secretKey].Secrets[secretKey], secretValue)
}
