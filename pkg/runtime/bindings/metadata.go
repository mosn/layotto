package bindings

type Metadata struct {
	Version  string
	Metadata map[string]string `json:"metadata"`
}
