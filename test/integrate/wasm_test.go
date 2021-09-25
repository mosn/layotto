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

package integrate

import (
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSayHello(t *testing.T) {
	ids := []string{"id_1", "id_2"}

	client := &http.Client{}
	req, _ := http.NewRequest("GET", "http://localhost:2045", nil)
	name := "Layotto"
	req.Header.Add("name", name)

	for _, id := range ids {
		req.Header.Set("id", id)
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Request failed, err: %s", err)
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Read body failed, err: %s", err)
		}
		assert.Equal(t, "Hi, "+name+"_"+id, string(body))
	}
}
