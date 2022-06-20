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
	"bytes"
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/jarcoal/httpmock"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/stretchr/testify/assert"

	"mosn.io/layotto/components/pkg/mock"
)

func TestNewClient(t *testing.T) {
	s := newQiniuOSSClient("1", "2", "3",
		"", true, false, false)

	assert.NotNil(t, s)
	assert.Equal(t, "1", s.AccessKey)
	assert.Equal(t, "2", s.SecretKey)
	assert.Equal(t, "3", s.Bucket)
}

func TestClientPut(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	m := mock.NewMockFormUploader(ctrl)
	bm := mock.NewMockBucketManager(ctrl)

	m.EXPECT().Put(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Eq("b/a.txt"), gomock.Any(), gomock.Any(), nil).Return(nil)
	m.EXPECT().Put(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Eq("a.txt"), gomock.Any(), gomock.Any(), nil).Return(nil)

	defer ctrl.Finish()

	s := newMockQiniuOSSClient("ak", "sk", "xc2022", "", true, m, bm)

	data := []byte("abc")
	err := s.put(context.Background(), "b/a.txt", bytes.NewReader(data), int64(len(data)))
	assert.NoError(t, err)

	err = s.put(context.Background(), "/a.txt", bytes.NewReader(data), int64(len(data)))
	assert.Error(t, err)

	err = s.put(context.Background(), "a.txt", bytes.NewReader(data), int64(len(data)))
	assert.NoError(t, err)

	err = s.put(context.Background(), "/b", bytes.NewReader(data), int64(len(data)))
	assert.Error(t, err)
}

func TestGetFromPrivate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	m := mock.NewMockFormUploader(ctrl)
	bm := mock.NewMockBucketManager(ctrl)
	defer ctrl.Finish()
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", `=~^http(s)?://example.com`,
		httpmock.NewStringResponder(200, ""))

	s := newMockQiniuOSSClient("ak", "sk", "xc2022", "https://example.com", true, m, bm)

	resp, err := s.get(context.Background(), "a.txt")
	assert.NoError(t, err)
	assert.NotNil(t, resp)

	resp, err = s.get(context.Background(), "/a.txt")
	assert.Error(t, err)
	assert.Nil(t, resp)

	s.Domain = "https://dsafafa.com"
	resp, err = s.get(context.Background(), "a.txt")
	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestGetFromPublic(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	m := mock.NewMockFormUploader(ctrl)
	bm := mock.NewMockBucketManager(ctrl)
	defer ctrl.Finish()
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", `=~^http(s)?://example.com`,
		httpmock.NewStringResponder(200, ""))

	s := newMockQiniuOSSClient("ak", "sk", "xc2022", "https://example.com", false, m, bm)

	resp, err := s.get(context.Background(), "a.txt")
	assert.NoError(t, err)
	assert.NotNil(t, resp)

	resp, err = s.get(context.Background(), "/a.txt")
	assert.Error(t, err)
	assert.Nil(t, resp)

	s.Domain = "https://dsafafa.com"
	resp, err = s.get(context.Background(), "a.txt")
	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestStatFile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	fu := mock.NewMockFormUploader(ctrl)
	bm := mock.NewMockBucketManager(ctrl)
	defer ctrl.Finish()

	bm.EXPECT().Stat(gomock.Eq("xc2022"), gomock.Eq("a.txt")).Return(storage.FileInfo{}, nil)
	s := newMockQiniuOSSClient("ak", "sk", "xc2022", "https://example.com", false, fu, bm)

	resp, err := s.stat(context.Background(), "a.txt")
	assert.NoError(t, err)
	assert.NotNil(t, resp)

	resp, err = s.stat(context.Background(), "/a.txt")
	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestDeleteFile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	fu := mock.NewMockFormUploader(ctrl)
	bm := mock.NewMockBucketManager(ctrl)
	defer ctrl.Finish()

	bm.EXPECT().Delete(gomock.Eq("xc2022"), gomock.Eq("a.txt")).Return(nil)
	bm.EXPECT().Delete(gomock.Eq("xc2022"), gomock.Eq("a.txt2")).Return(errors.New("aa"))
	s := newMockQiniuOSSClient("ak", "sk", "xc2022", "https://example.com", false, fu, bm)

	err := s.del(context.Background(), "a.txt")
	assert.NoError(t, err)

	err = s.del(context.Background(), "a.txt2")
	assert.Error(t, err)

	err = s.del(context.Background(), "/a.txt")
	assert.Error(t, err)
}

func TestListFile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	fu := mock.NewMockFormUploader(ctrl)
	bm := mock.NewMockBucketManager(ctrl)
	defer ctrl.Finish()

	bm.EXPECT().ListFiles(gomock.Eq("xc2022"), gomock.Eq("b/"), gomock.Any(), gomock.Any(), gomock.Any()).Return(make([]storage.ListItem, 0), make([]string, 0), "", false, nil)
	s := newMockQiniuOSSClient("ak", "sk", "xc2022", "https://example.com", false, fu, bm)

	items, prefixes, nextMarker, hasNext, err := s.list(context.Background(), "b/", 10, "")
	assert.NoError(t, err)
	assert.NotNil(t, items)
	assert.NotNil(t, prefixes)
	assert.Equal(t, "", nextMarker)
	assert.Equal(t, false, hasNext)

	items, prefixes, nextMarker, hasNext, err = s.list(context.Background(), "b/", 1001, "")
	assert.Error(t, err)
	assert.Nil(t, items)
	assert.Nil(t, prefixes)
	assert.Equal(t, "", nextMarker)
	assert.Equal(t, false, hasNext)

	items, prefixes, nextMarker, hasNext, err = s.list(context.Background(), "b/", -1, "")
	assert.Error(t, err)
	assert.Nil(t, items)
	assert.Nil(t, prefixes)
	assert.Equal(t, "", nextMarker)
	assert.Equal(t, false, hasNext)

}

func newMockQiniuOSSClient(ak, sk, bucket, domain string, private bool, fu FormUploader, bm BucketManager) *QiniuOSSClient {
	s := &QiniuOSSClient{
		AccessKey: ak,
		SecretKey: sk,
		Bucket:    bucket,
		fu:        fu,
		Domain:    domain,
		Private:   private,
		bm:        bm,
		mac:       qbox.NewMac("a", "b"),
	}

	return s
}
