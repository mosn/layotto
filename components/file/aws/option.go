package aws

import "github.com/jinzhu/copier"

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
)
