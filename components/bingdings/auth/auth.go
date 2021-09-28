package auth

import (
	"fmt"

	"github.com/dapr/components-contrib/bindings"
)

type AuthBindings struct {
}

// NewAuthBindings returns a new AuthBindings
func NewAuthBindings() *AuthBindings {
	return &AuthBindings{}
}

func (h *AuthBindings) Init(metadata bindings.Metadata) error {
	//do nothing
	return nil
}

func (h *AuthBindings) Invoke(req *bindings.InvokeRequest) (*bindings.InvokeResponse, error) {
	resp := &bindings.InvokeResponse{}
	// operation is request
	if req.Operation == "" {
		return nil, fmt.Errorf("illegal operation: %+v", req.Operation)
	}
	resp.Data = req.Data
	resp.Metadata = req.Metadata
	return resp, nil
}

func (h *AuthBindings) Operations() []bindings.OperationKind {
	return []bindings.OperationKind{
		bindings.CreateOperation, // For backward compatibility
		"get",
	}
}
