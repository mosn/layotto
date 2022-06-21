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

package aws

import (
	"testing"

	"github.com/jinzhu/copier"
)

func TestCopyWithConverterStrToStrPointer(t *testing.T) {
	type SrcStruct struct {
		Field1 string
	}

	type DestStruct struct {
		Field1 *string
	}

	src := SrcStruct{}

	var dst DestStruct

	err := copier.Copy(&src, &dst)

	if err != nil {
		t.Fatalf(`Should be able to copy from src to dst object. %v`, err)
		return
	}

	if src.Field1 != "" {
		t.Fatalf("got %q, wanted nil", src.Field1)
	}
}
