/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package sign

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"net/url"
	"strconv"
	"time"
)

const (
	httpHeaderAuthorization = "Authorization"
	httpHeaderTimestamp     = "Timestamp"

	authorizationFormat = "Apollo %s:%s"

	delimiter = "\n"
	question  = "?"
)

// AuthSignature apollo 授权
type AuthSignature struct {
}

// HTTPHeaders HTTPHeaders
func (t *AuthSignature) HTTPHeaders(url string, appID string, secret string) map[string][]string {
	ms := time.Now().UnixNano() / int64(time.Millisecond)
	timestamp := strconv.FormatInt(ms, 10)
	pathWithQuery := url2PathWithQuery(url)

	stringToSign := timestamp + delimiter + pathWithQuery
	signature := signString(stringToSign, secret)
	headers := make(map[string][]string, 2)

	signatures := make([]string, 0, 1)
	signatures = append(signatures, fmt.Sprintf(authorizationFormat, appID, signature))
	headers[httpHeaderAuthorization] = signatures

	timestamps := make([]string, 0, 1)
	timestamps = append(timestamps, timestamp)
	headers[httpHeaderTimestamp] = timestamps
	return headers
}

func signString(stringToSign string, accessKeySecret string) string {
	key := []byte(accessKeySecret)
	mac := hmac.New(sha1.New, key)
	mac.Write([]byte(stringToSign))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

func url2PathWithQuery(rawURL string) string {
	u, err := url.Parse(rawURL)
	if err != nil {
		return ""
	}
	pathWithQuery := u.Path

	if len(u.RawQuery) > 0 {
		pathWithQuery += question + u.RawQuery
	}
	return pathWithQuery
}
