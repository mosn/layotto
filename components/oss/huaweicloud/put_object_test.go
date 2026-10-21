/*
* Copyright 2026 Layotto Authors
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

package huaweicloud

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/huaweicloud/huaweicloud-sdk-go-obs/obs"
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

	client, err := obs.New("access-key", "secret-key", server.URL, obs.WithPathStyle(true), obs.WithSignature(obs.SignatureObs))
	require.NoError(t, err)
	defer client.Close()
	instance := &HuaweicloudOSS{client: client}

	_, putErr := instance.PutObject(context.Background(), &oss.PutObjectInput{
		Bucket:          "bucket",
		Key:             "object",
		DataStream:      strings.NewReader("hello"),
		ContentType:     "text/plain",
		ContentLanguage: "zh-CN",
		ContentMD5:      "XUFAKrxLKna5cZ2REBfFkg==",
	})
	require.NoError(t, putErr)
	assert.Equal(t, "hello", string(body))
	assert.Equal(t, "text/plain", headers.Get("Content-Type"))
	// This SDK version has no ContentLanguage field in PutObjectInput.
	assert.Empty(t, headers.Get("Content-Language"))
	assert.Equal(t, "XUFAKrxLKna5cZ2REBfFkg==", headers.Get("Content-MD5"))

}
