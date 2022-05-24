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
	"context"
	"io"
)

type Oss interface {
	InitConfig(context.Context, *FileConfig) error
	InitClient(context.Context, *InitRequest) error
	GetObject(context.Context, *GetObjectInput) (io.ReadCloser, error)
	PutObject(context.Context) error
}

type BaseConfig struct {
}
type InitRequest struct {
	App      string
	Metadata map[string]string
}

type GetObjectInput struct {
	ClientName                 string
	Bucket                     string
	ExpectedBucketOwner        string
	IfMatch                    string
	IfModifiedSince            string
	IfNoneMatch                string
	IfUnmodifiedSince          string
	Key                        string
	PartNumber                 int64
	Range                      string
	RequestPayer               string
	ResponseCacheControl       string
	ResponseContentDisposition string
	ResponseContentEncoding    string
	ResponseContentLanguage    string
	ResponseContentType        string
	ResponseExpires            string
	SseCustomerAlgorithm       string
	SseCustomerKey             string
	SseCustomerKeyMd5          string
	VersionId                  string
}

type PutObjectInput struct {
	ClientName       string
	Bucket           string
	Key              string
	DataStream       io.Reader
	Acl              string
	BucketKeyEnabled bool
	CacheControl     string
}
