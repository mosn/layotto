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
