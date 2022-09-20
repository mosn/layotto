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
package snowflake

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPaddingRingBuffer(t *testing.T) {
	rb := NewRingBuffer(8192)

	rb.TimeBits = 28
	rb.WorkIdBits = 22
	rb.SeqBits = 13
	rb.PaddingFactor = 50

	s := "2022-01-01"

	var tmp time.Time
	tmp, err := time.ParseInLocation("2006-01-02", s, time.Local)
	assert.NoError(t, err)
	startTime := tmp.Unix()
	rb.CurrentTimeStamp = time.Now().Unix() - startTime
	rb.PaddingRingBuffer()

	var uid int64
	for i := 0; i < 8192; i++ {
		uid, err = rb.Take()
		assert.NoError(t, err)
		assert.NotEqual(t, uid, 0)
	}
	uid, err = rb.Take()
	assert.NoError(t, err)
	assert.NotEqual(t, uid, 0)
}
