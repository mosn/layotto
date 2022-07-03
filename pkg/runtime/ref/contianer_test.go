package ref

import (
	"github.com/stretchr/testify/assert"
	mock "mosn.io/layotto/pkg/mock"
	"mosn.io/layotto/pkg/mock/components/secret"
	"testing"
)

func TestRefContainer(t *testing.T) {

	container := NewRefContainer()

	ss := &secret.FakeSecretStore{}
	container.SecretRef["fake"] = ss
	cf := &mock.MockStore{}
	container.ConfigRef["mock"] = cf
	assert.Equal(t, ss, container.GetSecretStore("fake"))
	assert.Equal(t, cf, container.GetConfigStore("mock"))

}
