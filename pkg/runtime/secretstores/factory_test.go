package secretstores

import (
	secretstores_loader "github.com/dapr/components-contrib/secretstores"
	"mosn.io/layotto/pkg/mock/components/secret"
	"testing"
)

func TestNewFactory(t *testing.T) {
	NewFactory("test", func() secretstores_loader.SecretStore {
		return secret.FakeSecretStore{}
	})
}
