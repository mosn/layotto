package common

import (
	"mosn.io/layotto/components/configstores"
)

type InjectComponent interface {
	InjectConfigComponent(configs []configstores.Store) (err error)
}
