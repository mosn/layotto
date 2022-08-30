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

package ceph

import (
	"time"

	"github.com/jinzhu/copier"
)

var (
	int642time = copier.TypeConverter{
		SrcType: int64(0),
		DstType: &time.Time{},
		Fn: func(src interface{}) (interface{}, error) {
			s, _ := src.(int64)
			t := time.Unix(s, 0)
			return &t, nil
		},
	}
	time2int64 = copier.TypeConverter{
		SrcType: &time.Time{},
		DstType: int64(0),
		Fn: func(src interface{}) (interface{}, error) {
			s, _ := src.(*time.Time)
			return s.Unix(), nil
		},
	}
)
