package aws

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
	srcValue := &ValueWithInt64{TestInt64toTime: timer, TestString: ""}
	destValue := &ValueWithTimer{}
	err := copier.CopyWithOption(destValue, srcValue, copier.Option{IgnoreEmpty: true, DeepCopy: true, Converters: []copier.TypeConverter{str2point, int642time}})
	assert.Nil(t, err)
	assert.Nil(t, destValue.TestString)
	assert.Equal(t, timer, destValue.TestInt64toTime.Unix())

	ti := time.Now()
	src := &ValueWithTimer{TestInt64toTime: &ti}
	dst := &ValueWithInt64{}
	err = copier.CopyWithOption(dst, src, copier.Option{IgnoreEmpty: true, DeepCopy: true, Converters: []copier.TypeConverter{time2int64}})
	assert.Nil(t, err)
	assert.Equal(t, ti.Unix(), dst.TestInt64toTime)
}
