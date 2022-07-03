package ref

import (
	"github.com/stretchr/testify/assert"
	"mosn.io/layotto/components/ref"
	"mosn.io/layotto/pkg/mock"
	"mosn.io/layotto/pkg/mock/components/secret"
	"testing"
)

func TestInject(t *testing.T) {

	container := NewRefContainer()

	ss := &secret.FakeSecretStore{}
	container.SecretRef["fake_secret_store"] = ss
	cf := &mock.MockStore{}
	container.ConfigRef["mock_config_store"] = cf

	injector := DefaultInjector{Container: *container}
	meta := make(map[string]string)
	var items []*ref.Item
	items = append(items, &ref.Item{
		Name: "fake_secret_store",
		Key:  "good-key",
	})
	injector.InjectSecretRef(items, meta)
	assert.Equal(t, meta["good-key"], "life is good")

}
