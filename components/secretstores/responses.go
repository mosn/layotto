package secretstores

// GetSecretResponse describes the response object for a secret returned from a secret store.
type GetSecretResponse struct {
	Data map[string]string `json:"data"`
}

// BulkGetSecretResponse describes the response object for all the secrets returned from a secret store.
type BulkGetSecretResponse struct {
	Data map[string]map[string]string `json:"data"`
}
