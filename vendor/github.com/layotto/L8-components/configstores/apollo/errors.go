package apollo

import (
	"errors"
	"fmt"
)

var ErrNoConfig = errors.New("configuration illegal:no config data")

func errConfigMissingField(field string) error {
	return errors.New(fmt.Sprintf("configuration illegal:no %s", field))
}

func errParamsMissingField(field string) error {
	return errors.New(fmt.Sprintf("params illegal:no %s", field))
}
