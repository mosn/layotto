package state

// Config wraps configuration for a state implementation
type Config struct {
	Metadata map[string]string `json:"metadata"`
}
