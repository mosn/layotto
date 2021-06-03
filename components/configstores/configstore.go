package configstores

import (
	"context"
)

const (
	ServiceName = "configStore"
	All         = "*"
)

const (
	AppId = iota
	Group
	Label
	Key
	Tag
)

// Store is an interface to perform operations on config store
type Store interface {
	//Init init the configuration store.
	Init(config *StoreConfig) error

	// GetSpecificKeysValue get specific key value.
	Get(context.Context, *GetRequest) ([]*ConfigurationItem, error)

	// Set saves configuration into configuration store.
	Set(context.Context, *SetRequest) error

	// Delete deletes configuration from configuration store.
	Delete(context.Context, *DeleteRequest) error

	// Subscribe subscribe the configurations updates.
	Subscribe(*SubscribeReq, chan *SubscribeResp) error

	//StopSubscribe stop subs
	StopSubscribe()

	GetDefaultGroup() string

	GetDefaultLabel() string
}
