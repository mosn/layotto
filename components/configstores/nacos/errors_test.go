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

package nacos

import "testing"

func Test_errConfigMissingField(t *testing.T) {
	err := errConfigMissingField("")
	if s := err.Error(); s == "" {
		t.Errorf("error has no text")
	}
}

func Test_errParamsMissingField(t *testing.T) {
	err := errParamsMissingField("")
	if s := err.Error(); s == "" {
		t.Errorf("error has no text")
	}
}
