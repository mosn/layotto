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
// CODE ATTRIBUTION: https://github.com/dapr/dapr
// We copied these code here to make our runtime compatible with dapr's component.
package state

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

const (
	strategyKey = "keyPrefix"

	strategyAppid     = "appid"
	strategyStoreName = "name"
	strategyNone      = "none"
	strategyDefault   = strategyNone

	daprSeparator = "||"
)

var statesConfiguration = map[string]*StoreConfiguration{}

type StoreConfiguration struct {
	keyPrefixStrategy string
}

// Save StateConfiguration by storeName
func SaveStateConfiguration(storeName string, metadata map[string]string) error {
	// convert
	strategy := metadata[strategyKey]
	// Change strategy to lowercase
	strategy = strings.ToLower(strategy)
	//if strategy is "",use default values("none")
	if strategy == "" {
		strategy = strategyDefault
	} else {
		// Check if the secret key is legitimate
		err := checkKeyIllegal(metadata[strategyKey])
		if err != nil {
			return err
		}
	}
	// convert
	statesConfiguration[storeName] = &StoreConfiguration{keyPrefixStrategy: strategy}
	return nil
}

func GetModifiedStateKey(key, storeName, appID string) (string, error) {
	// Check if the secret key is legitimate
	if err := checkKeyIllegal(key); err != nil {
		return "", err
	}
	// Get stateConfiguration by storeName
	stateConfiguration := getStateConfiguration(storeName)
	// Determine the keyPrefixStrategy type
	switch stateConfiguration.keyPrefixStrategy {
	case strategyNone:
		return key, nil
	case strategyStoreName:
		return fmt.Sprintf("%s%s%s", storeName, daprSeparator, key), nil
	case strategyAppid:
		if appID == "" {
			return key, nil
		}
		return fmt.Sprintf("%s%s%s", appID, daprSeparator, key), nil
	default:
		return fmt.Sprintf("%s%s%s", stateConfiguration.keyPrefixStrategy, daprSeparator, key), nil
	}
}

func GetOriginalStateKey(modifiedStateKey string) string {
	// Split modifiedStateKey by daprSeparator("||")
	splits := strings.Split(modifiedStateKey, daprSeparator)
	if len(splits) <= 1 {
		return modifiedStateKey
	}
	return splits[1]
}

func getStateConfiguration(storeName string) *StoreConfiguration {
	// Get statesConfiguration by storeName
	c := statesConfiguration[storeName]
	// If statesConfiguration is empty, strategyDefault("none") is provided
	if c == nil {
		c = &StoreConfiguration{keyPrefixStrategy: strategyDefault}
		statesConfiguration[storeName] = c
	}

	return c
}

func checkKeyIllegal(key string) error {
	// Determine if the key contains daprSeparator
	if strings.Contains(key, daprSeparator) {
		return errors.Errorf("input key/keyPrefix '%s' can't contain '%s'", key, daprSeparator)
	}
	return nil
}
