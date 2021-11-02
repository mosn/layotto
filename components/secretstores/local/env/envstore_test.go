package env

import (
	"github.com/stretchr/testify/assert"
	"mosn.io/layotto/components/secretstores"
	"mosn.io/pkg/log"
	"os"
	"testing"
)

func TestInit(t *testing.T) {
	secret := "secret1"
	key := "TEST_SECRET"

	s := envSecretStore{logger: log.Logger{}}

	os.Setenv(key, secret)
	assert.Equal(t, secret, os.Getenv(key))

	t.Run("Test init", func(t *testing.T) {
		err := s.Init(secretstores.Metadata{})
		assert.Nil(t, err)
	})

	t.Run("Test set and get", func(t *testing.T) {
		err := s.Init(secretstores.Metadata{})
		assert.Nil(t, err)
		resp, err := s.GetSecret(secretstores.GetSecretRequest{Name: key})
		assert.Nil(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, secret, resp.Data[key])
	})

	t.Run("Test bulk get", func(t *testing.T) {
		err := s.Init(secretstores.Metadata{})
		assert.Nil(t, err)
		resp, err := s.BulkGetSecret(secretstores.BulkGetSecretRequest{})
		assert.Nil(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, secret, resp.Data[key][key])
	})
}
