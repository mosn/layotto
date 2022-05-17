package ref

import (
	"context"
	"github.com/dapr/components-contrib/secretstores"
	"mosn.io/layotto/components/configstores"
)

const (
	TypeConfig = "config_store"
	TypeSecret = "secret_store"
)

type DefaultInjector struct {
	Container RefContainer
}

//InjectSecretCmpRef inject secret Component
func (i *DefaultInjector) InjectSecretCmpRef(componentRefs []*ComponentRef) map[string]secretstores.SecretStore {
	cmpRefs := make(map[string]secretstores.SecretStore)
	for _, ref := range componentRefs {
		if ref.Type != TypeSecret {
			continue
		}
		store := i.Container.GetSecretStore(ref.Name)
		if store == nil {
			continue
		}
		cmpRefs[ref.Name] = store
	}
	return cmpRefs
}

func (i *DefaultInjector) InjectConfigCmpRef(componentRefs []*ComponentRef) map[string]configstores.Store {
	cmpRefs := make(map[string]configstores.Store)
	for _, ref := range componentRefs {
		if ref.Type != TypeConfig {
			continue
		}
		store := i.Container.GetConfigStore(ref.Name)
		if store == nil {
			continue
		}
		cmpRefs[ref.Name] = store
	}
	return cmpRefs
}

//InjectSecretRef  inject secret to metaData
func (i *DefaultInjector) InjectSecretRef(items []*RefItem, metaData map[string]string) error {

	meta := make(map[string]string)
	for _, item := range items {
		store := i.Container.GetSecretStore(item.Name)
		secret, err := store.GetSecret(secretstores.GetSecretRequest{
			Name: item.Key,
		})
		if err != nil {
			return err
		}
		for k, v := range secret.Data {
			meta[k] = v
		}
	}
	//avoid part of assign because of err
	for k, v := range meta {
		metaData[k] = v
	}
	return nil
}

//InjectConfigRef inject config
func (i *DefaultInjector) InjectConfigRef(items []*RefItem, metaData map[string]string) error {
	//TODO: how to inject config  subscribe or once ?  It can be configured by the user
	meta := make(map[string]string)
	for _, item := range items {
		store := i.Container.GetConfigStore(item.Name)
		cfgS, err := store.Get(context.Background(), &configstores.GetRequest{
			Keys: []string{item.Key},
		})
		if err != nil {
			return err
		}

		//TODO： we can put the code of subscribe here

		for _, cfg := range cfgS {
			meta[cfg.Key] = cfg.Content
		}
	}
	for k, v := range meta {
		metaData[k] = v
	}
	return nil
}
