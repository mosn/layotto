package file

import (
	"encoding/json"
	"io"
)

// FileConfig wraps configuration for a file implementation
type FileConfig struct {
	Metadata json.RawMessage
}

type PutFileStu struct {
	DataStream io.Reader
	FileName   string
	Metadata   map[string]string
}

type GetFileStu struct {
	FileName string
	Metadata map[string]string
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
