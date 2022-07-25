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

package qiniu

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
)

type QiniuOSSClient struct {
	AccessKey string
	SecretKey string
	Bucket    string
	Domain    string
	Private   bool

	mac *qbox.Mac
	fu  FormUploader
	bm  BucketManager
}

type FormUploader interface {
	Put(ctx context.Context, ret interface{}, uptoken, key string, data io.Reader, size int64, extra *storage.PutExtra) (err error)
}

type BucketManager interface {
	Stat(bucket, key string) (storage.FileInfo, error)
	Delete(bucket, key string) (err error)
	ListFiles(bucket, prefix, delimiter, marker string,
		limit int) (entries []storage.ListItem, commonPrefixes []string, nextMarker string, hasNext bool, err error)
}

func newQiniuOSSClient(ak, sk, bucket, domain string, private bool, useHttps, userCdnDomains bool) *QiniuOSSClient {
	cfg := storage.Config{
		UseHTTPS:      useHttps,
		UseCdnDomains: userCdnDomains,
	}

	mac := qbox.NewMac(ak, sk)
	s := &QiniuOSSClient{
		AccessKey: ak,
		SecretKey: sk,
		Bucket:    bucket,
		fu:        storage.NewFormUploader(&cfg),
		Domain:    domain,
		Private:   private,
		mac:       mac,
		bm:        storage.NewBucketManager(mac, &cfg),
	}
	return s
}

func (s *QiniuOSSClient) put(ctx context.Context, fileName string, data io.Reader, dataSize int64) error {
	if err := s.checkFileName(fileName); err != nil {
		return err
	}

	putPolicy := storage.PutPolicy{
		Scope: s.Bucket,
	}

	upToken := putPolicy.UploadToken(qbox.NewMac(s.AccessKey, s.SecretKey))

	ret := storage.PutRet{}
	err := s.fu.Put(ctx, &ret, upToken, fileName, data, dataSize, nil)

	return err
}

func (s *QiniuOSSClient) get(_ context.Context, fileName string) (io.ReadCloser, error) {
	if err := s.checkFileName(fileName); err != nil {
		return nil, err
	}

	var accessUrl string

	if !s.Private {
		accessUrl = storage.MakePublicURL(s.Domain, fileName)

	} else {
		deadline := time.Now().Add(time.Second * 60).Unix() //1小时有效期
		accessUrl = storage.MakePrivateURL(s.mac, s.Domain, fileName, deadline)
	}

	resp, err := http.Get(accessUrl)
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}

func (s *QiniuOSSClient) stat(_ context.Context, fileName string) (*storage.FileInfo, error) {
	if err := s.checkFileName(fileName); err != nil {
		return nil, err
	}

	resp, err := s.bm.Stat(s.Bucket, fileName)
	return &resp, err
}

func (s *QiniuOSSClient) del(_ context.Context, fileName string) error {
	if err := s.checkFileName(fileName); err != nil {
		return err
	}

	return s.bm.Delete(s.Bucket, fileName)
}

func (s *QiniuOSSClient) list(_ context.Context, prefix string, limit int, marker string) (entries []storage.ListItem, commonPrefixes []string, nextMarker string, hasNext bool, err error) {
	if limit > 1000 {
		return nil, nil, "", false, errors.New("limit must be <=1000")
	}

	if limit <= 0 {
		return nil, nil, "", false, errors.New("limit must be >0")
	}

	return s.bm.ListFiles(s.Bucket, prefix, "", marker, limit)
}

func (s *QiniuOSSClient) checkFileName(fileName string) error {
	index := strings.Index(fileName, "/")
	if index == 0 {
		return fmt.Errorf("invalid fileName format")
	}

	return nil
}
