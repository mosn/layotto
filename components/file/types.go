package file

// FileConfig wraps configuration for a file implementation
type FileConfig struct {
	Metadata map[string]string `json:"metadata"`
}

type PutFileStu struct {
	Data        []byte
	FileName    string
	Metadata    map[string]string
	StreamId    int64
	ChunkNumber int
}

type GetFileStu struct {
	ObjectName string
	Metadata   map[string]string
}
