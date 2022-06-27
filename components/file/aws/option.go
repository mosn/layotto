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
	"errors"
	"time"

	"github.com/jinzhu/copier"
)

var str = "'"
var (
	str2point = copier.TypeConverter{
		SrcType: copier.String,
		DstType: &str,
		Fn: func(src interface{}) (interface{}, error) {
			s, _ := src.(string)
			// return nil on empty string
			if s == "" {
				return nil, nil
			}
			return &s, nil
		},
	}
	int642time = copier.TypeConverter{
		SrcType: int64(0),
		DstType: &time.Time{},
		Fn: func(src interface{}) (interface{}, error) {
			s, ok := src.(int64)
			if !ok {
				return nil, errors.New("src type not matching")
			}
			t := time.Unix(0, s)
			return &t, nil
		},
	}
	time2int64 = copier.TypeConverter{
		SrcType: &time.Time{},
		DstType: int64(0),
		Fn: func(src interface{}) (interface{}, error) {
			s, ok := src.(*time.Time)
			if !ok {
				return nil, errors.New("src type not matching")
			}
			return s.Unix(), nil
		},
	}
)
