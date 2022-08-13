package common

import "context"

type DynamicComponent interface {
	ApplyConfig(ctx context.Context, metadata map[string]string) (err error)
}
