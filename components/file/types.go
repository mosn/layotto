package file

// FileConfig wraps configuration for a file implementation
type FileConfig struct {
	Metadata []map[string]interface{} `json:"metadata"`
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

type DelRequest struct {
	FileName string
	Metadata map[string]string
}

type ListRequest struct {
	DirectoryName string
	Metadata      map[string]string
}

type ListResp struct {
	FilesName []string
}
