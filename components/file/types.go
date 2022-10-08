/*
 * Copyright 2021 Layotto Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package file

import (
	"encoding/json"
	"io"

	"mosn.io/layotto/components/ref"
)

// FileConfig wraps configuration for a file implementation
type FileConfig struct {
	ref.Config
	Metadata json.RawMessage `json:"metadata"`
	Type     string          `json:"type"`
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
	Marker        string
	PageSize      int32
	Metadata      map[string]string
}

type FilesInfo struct {
	FileName     string
	Size         int64
	LastModified string
	Meta         map[string]string
}

type ListResp struct {
	Files       []*FilesInfo
	Marker      string
	IsTruncated bool
}

type FileMetaRequest struct {
	FileName string
	Metadata map[string]string
}

type FileMetaResp struct {
	Size         int64
	LastModified string
	Metadata     map[string][]string
}
