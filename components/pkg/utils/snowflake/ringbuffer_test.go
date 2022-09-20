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
