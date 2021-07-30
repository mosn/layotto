package sequencer

import (
	"fmt"
	"github.com/pkg/errors"
	"strings"
)

const (
	separator    = "|||"
	commonPrefix = "sequencer"
)

func GetModifiedKey(key, storeName, appID string) (string, error) {
	if err := checkKeyIllegal(key); err != nil {
		return "", err
	}
	return fmt.Sprintf("%s%s%s", commonPrefix, separator, key), nil
}

func checkKeyIllegal(key string) error {
	if strings.Contains(key, separator) {
		return errors.Errorf("input key/keyPrefix '%s' can't contain '%s'", key, separator)
	}
	return nil
}
