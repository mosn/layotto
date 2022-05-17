package ref

import (
	"github.com/dapr/components-contrib/secretstores"
	"mosn.io/layotto/components/configstores"
)

// RefContainer  hold all secret&config store
type RefContainer struct {
	SecretRef map[string]secretstores.SecretStore
	ConfigRef map[string]configstores.Store
}

func NewRefContainer() *RefContainer {
	return &RefContainer{
		SecretRef: make(map[string]secretstores.SecretStore),
		ConfigRef: make(map[string]configstores.Store),
	}
}

func (r *RefContainer) GetSecretStore(key string) secretstores.SecretStore {
	return r.SecretRef[key]
}

func (r *RefContainer) GetConfigStore(key string) configstores.Store {
	return r.ConfigRef[key]
}
