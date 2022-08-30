package runtime

import "mosn.io/layotto/components/lock"

type GeneratedConfig struct {
	LockManagement map[string]lock.Config `json:"lock"`
}
