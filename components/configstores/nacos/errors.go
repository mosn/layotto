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

import (
	"errors"
	"fmt"
)

var (
	// IllegalParam
	// This error is usually caused by the logic of the parameters,
	// such as filling in parameters with illegal characters
	IllegalParam = errors.New("illegal parameter")
)

func errConfigMissingField(field string) error {
	return fmt.Errorf("configuration illegal:no %s", field)
}

func errParamsMissingField(field string) error {
	return fmt.Errorf("params illegal:no %s", field)
}
