//
// Copyright 2021 Layotto Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package aliyun

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	aliyunSDK "github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"mosn.io/layotto/components/oss"
)

func TestPutObjectContentHeaders(t *testing.T) {
	var headers http.Header
	var body []byte
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPut, r.Method)
		headers = r.Header.Clone()
		var err error
		body, err = io.ReadAll(r.Body)
		assert.NoError(t, err)
		w.Header().Set("ETag", "\"5d41402abc4b2a76b9719d911017c592\"")
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client, err := aliyunSDK.New(server.URL, "access-key", "secret-key")
	require.NoError(t, err)
	instance := &AliyunOSS{client: client}
	for _, tc := range []struct {
		name               string
		signedURL          string
		withContentHeaders bool
		wantContentType    string
	}{
		{name: "normal", withContentHeaders: true, wantContentType: "text/plain"},
		{name: "signed URL", signedURL: server.URL + "/bucket/object", withContentHeaders: true, wantContentType: "text/plain"},
		// Without ContentType, the SDK supplies a default only for normal uploads.
		{name: "normal without content headers", wantContentType: "application/octet-stream"},
		{name: "signed URL without content headers", signedURL: server.URL + "/bucket/object"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			headers = nil
			body = nil
			req := &oss.PutObjectInput{
				Bucket:     "bucket",
				Key:        "object",
				DataStream: strings.NewReader("hello"),
				SignedUrl:  tc.signedURL,
			}
			if tc.withContentHeaders {
				req.ContentType = "text/plain"
				req.ContentLanguage = "zh-CN"
				req.ContentMD5 = "XUFAKrxLKna5cZ2REBfFkg=="
			}
			_, putErr := instance.PutObject(context.Background(), req)
			require.NoError(t, putErr)
			assert.Equal(t, "hello", string(body))
			assert.Equal(t, tc.wantContentType, headers.Get("Content-Type"))
			assert.Equal(t, req.ContentLanguage, headers.Get("Content-Language"))
			assert.Equal(t, req.ContentMD5, headers.Get("Content-MD5"))
			if !tc.withContentHeaders {
				assert.NotContains(t, headers, "Content-Language")
				assert.NotContains(t, headers, "Content-Md5")
				if tc.signedURL != "" {
					assert.NotContains(t, headers, "Content-Type")
				}
			}
		})
	}
}
