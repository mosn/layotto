package secretstores

import (
	"testing"

	"mosn.io/layotto/pkg/mock/components/secret"

	"github.com/dapr/components-contrib/secretstores"
)

func TestNewFactory(t *testing.T) {
	NewFactory("test", func() secretstores.SecretStore {
		return secret.FakeSecretStore{}
	})
}
