package secretstores

import (
	"github.com/dapr/components-contrib/secretstores"
	"mosn.io/layotto/pkg/mock/components/secret"
	"testing"
)

func TestNewFactory(t *testing.T) {
	NewFactory("test", func() secretstores.SecretStore {
		return secret.FakeSecretStore{}
	})
}
