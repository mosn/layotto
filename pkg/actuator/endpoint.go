package actuator

import (
	"context"
)

type Endpoint interface {
	Handle(ctx context.Context, params ParamsScanner) (jsonObject map[string]interface{}, err error)
}

type ParamsScanner interface {
	//Next get the next param.
	Next() string
	HasNext() bool
}
