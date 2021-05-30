package pubsub

// StoreConfig wraps configuration for a store implementation
type Config struct {
	AppId            string            `json:"app_id"`
	GrpcCallbackPort int               `json:"grpc_callback_port"`
	Metadata         map[string]string `json:"metadata"`
}
