package secret

import (
	"context"
	"errors"

	"github.com/dapr/components-contrib/common/features"
	"github.com/dapr/components-contrib/secretstores"
)

type FakeSecretStore struct{}

type Feature = features.Feature[secretstores.SecretStore]

func (c FakeSecretStore) GetSecret(ctx context.Context, req secretstores.GetSecretRequest) (secretstores.GetSecretResponse, error) {
	if req.Name == "good-key" {
		return secretstores.GetSecretResponse{
			Data: map[string]string{"good-key": "life is good"},
		}, nil
	}

	if req.Name == "error-key" {
		return secretstores.GetSecretResponse{}, errors.New("error occurs with error-key")
	}

	return secretstores.GetSecretResponse{}, nil
}

func (c FakeSecretStore) BulkGetSecret(ctx context.Context, req secretstores.BulkGetSecretRequest) (secretstores.BulkGetSecretResponse, error) {
	response := map[string]map[string]string{}

	response["good-key"] = map[string]string{"good-key": "life is good"}

	return secretstores.BulkGetSecretResponse{
		Data: response,
	}, nil
}

func (c FakeSecretStore) Init(ctx context.Context, metadata secretstores.Metadata) error {
	return nil
}

func (c FakeSecretStore) Close() error {
	return nil
}

func (c FakeSecretStore) Features() []Feature {
	return nil
}
