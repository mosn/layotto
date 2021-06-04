package client

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSaveConfiguration(t *testing.T) {
	ctx := context.Background()
	item1 := &ConfigurationItem{Key: "hello1", Content: "world1"}
	item2 := &ConfigurationItem{Key: "hello2", Content: "world2"}
	saveRequest := &SaveConfigurationRequest{StoreName: "etcd", AppId: "sofa"}
	saveRequest.Items = append(saveRequest.Items, item1)
	saveRequest.Items = append(saveRequest.Items, item2)
	t.Run("save configuration", func(t *testing.T) {
		err := testClient.SaveConfiguration(ctx, saveRequest)
		assert.Nil(t, err)
	})
}
func TestGetConfiguration(t *testing.T) {
	getRequest := &ConfigurationRequestItem{StoreName: "etcd", AppId: "sofa", Keys: []string{"hello1", "hello2"}}
	t.Run("get configuration", func(t *testing.T) {
		resp, err := testClient.GetConfiguration(context.Background(), getRequest)
		assert.Nil(t, err)
		assert.Equal(t, resp[0].Key, "hello1")
		assert.Equal(t, resp[0].Content, "world1")
		assert.Equal(t, resp[1].Key, "hello2")
		assert.Equal(t, resp[1].Content, "world2")
	})
}

func TestDeleteConfiguration(t *testing.T) {
	ctx := context.Background()
	deleteRequest := &ConfigurationRequestItem{StoreName: "etcd", AppId: "sofa", Keys: []string{"hello1", "hello2"}}
	t.Run("delete configuration", func(t *testing.T) {
		err := testClient.DeleteConfiguration(ctx, deleteRequest)
		assert.Nil(t, err)
	})

	getRequest := &ConfigurationRequestItem{StoreName: "etcd", AppId: "sofa", Keys: []string{"hello1", "hello2"}}
	t.Run("get configuration", func(t *testing.T) {
		resp, err := testClient.GetConfiguration(context.Background(), getRequest)
		assert.Nil(t, err)
		assert.Equal(t, 0, len(resp))
	})
}

func TestSubscribeConfiguration(t *testing.T) {
	item := &ConfigurationRequestItem{StoreName: "etcd", AppId: "sofa", Keys: []string{"hello1"}}
	ch := testClient.SubscribeConfiguration(context.Background(), item)
	for wc := range ch {
		assert.Equal(t, wc.Item.Items[0].Key, "hello1")
		assert.Equal(t, wc.Item.Items[0].Content, "Test")
	}
}
