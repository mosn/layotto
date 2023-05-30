package nacos

import (
	"fmt"
)

func errConfigMissingField(field string) error {
	return fmt.Errorf("configuration illegal:no %s", field)
}

func errParamsMissingField(field string) error {
	return fmt.Errorf("params illegal:no %s", field)
}
