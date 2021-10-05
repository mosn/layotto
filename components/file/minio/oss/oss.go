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

package oss

import (
	"io"

	"github.com/minio/minio-go"
	"github.com/minio/minio-go/v6"
	"mosn.io/layotto/components/file"
)

type MinioOss struct {
	client map[string]*minio.Client
	meta   map[string]*MinioMetaData
}

type MinioMetaData struct {
	EndPoint        string `json:"endpoint"`
	AccessKeyID     string `json:"accessKeyID"`
	AccessKeySecret string `json:"accessKeySecret"`
	UseSSL          bool   `json:"useSSL`
}

func (m *MinioOss) Init(config *file.FileConfig) error {
	return nil
}

func (m *MinioOss) Put(st *file.PutFileStu) error {
	return nil
}

func (m *MinioOss) Get(st *file.GetFileStu) (io.ReadCloser, error) {
	return nil, nil
}

func (m *MinioOss) List(st *file.ListRequest) (*ListResp, error) {
	return nil, nil
}

func (m *MinioOss) Del(st *file.DelRequest) error {
	return nil
}
