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

package oss

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/jinzhu/copier"
)

func TestCopierOption(t *testing.T) {
	type ValueWithInt64 struct {
		TestString      string
		TestInt64toTime int64
	}

	type ValueWithTimer struct {
		TestString      *string
		TestInt64toTime *time.Time
	}
	timer := time.Now().Unix()
	srcValue := &ValueWithInt64{TestInt64toTime: timer}
	destValue := &ValueWithTimer{}
	err := copier.CopyWithOption(destValue, srcValue, copier.Option{IgnoreEmpty: true, DeepCopy: true, Converters: []copier.TypeConverter{Int64ToTime}})
	assert.Nil(t, err)
	assert.Nil(t, destValue.TestString)
	assert.Equal(t, timer, destValue.TestInt64toTime.Unix())

	ti := time.Now()
	src := &ValueWithTimer{TestInt64toTime: &ti}
	dst := &ValueWithInt64{}
	err = copier.CopyWithOption(dst, src, copier.Option{IgnoreEmpty: true, DeepCopy: true, Converters: []copier.TypeConverter{TimeToInt64}})
	assert.Nil(t, err)
	assert.Equal(t, ti.Unix(), dst.TestInt64toTime)
}
