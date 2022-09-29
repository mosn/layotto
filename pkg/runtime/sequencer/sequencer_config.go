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
package sequencer

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

const (
	strategyKey       = "keyPrefix"
	strategyAppid     = "appid"
	strategyStoreName = "name"
	strategyNone      = "none"
	strategyDefault   = strategyAppid
	apiPrefix         = "sequencer"
	apiSeparator      = "|||"
	separator         = "||"
)

var seqConfiguration = map[string]*StoreConfiguration{}

type StoreConfiguration struct {
	keyPrefixStrategy string
}

func SaveSeqConfiguration(storeName string, metadata map[string]string) error {
	strategy := strings.ToLower(metadata[strategyKey])
	if strategy == "" {
		strategy = strategyDefault
	} else {
		err := checkKeyIllegal(metadata[strategyKey])
		if err != nil {
			return err
		}
	}

	seqConfiguration[storeName] = &StoreConfiguration{keyPrefixStrategy: strategy}
	return nil
}

func GetModifiedSeqKey(key, storeName, appID string) (string, error) {
	if err := checkKeyIllegal(key); err != nil {
		return "", err
	}
	config := getConfiguration(storeName)
	switch config.keyPrefixStrategy {
	case strategyNone:
		return fmt.Sprintf("%s%s%s", apiPrefix, apiSeparator, key), nil
	case strategyStoreName:
		return fmt.Sprintf("%s%s%s%s%s", apiPrefix, apiSeparator, storeName, separator, key), nil
	case strategyAppid:
		if appID == "" {
			return fmt.Sprintf("%s%s%s", apiPrefix, apiSeparator, key), nil
		}
		return fmt.Sprintf("%s%s%s%s%s", apiPrefix, apiSeparator, appID, separator, key), nil
	default:
		return fmt.Sprintf("%s%s%s%s%s", apiPrefix, apiSeparator, config.keyPrefixStrategy, separator, key), nil
	}
}

func getConfiguration(storeName string) *StoreConfiguration {
	c := seqConfiguration[storeName]
	if c == nil {
		c = &StoreConfiguration{keyPrefixStrategy: strategyDefault}
		seqConfiguration[storeName] = c
	}

	return c
}

func checkKeyIllegal(key string) error {
	if strings.Contains(key, separator) {
		return errors.Errorf("input key/keyPrefix '%s' can't contain '%s'", key, separator)
	}
	return nil
}
