package pubsub

// Config wraps configuration for a pubsub implementation
type Config struct {
	Metadata map[string]string `json:"metadata"`
}
