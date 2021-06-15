// CODE ATTRIBUTION: https://github.com/dapr/dapr
// We copied these code here to make our runtime compatible with dapr's component.
package state

import (
	"fmt"
	runtimev1pb "mosn.io/layotto/spec/proto/runtime/v1"
	"strings"

	"github.com/pkg/errors"
)

const (
	strategyKey = "keyPrefix"

	strategyAppid     = "appid"
	strategyStoreName = "name"
	strategyNone      = "none"
	strategyDefault   = strategyAppid

	daprSeparator = "||"
)

var statesConfiguration = map[string]*StoreConfiguration{}

type StoreConfiguration struct {
	keyPrefixStrategy string
}

func SaveStateConfiguration(storeName string, metadata map[string]string) error {
	strategy := metadata[strategyKey]
	strategy = strings.ToLower(strategy)
	if strategy == "" {
		strategy = strategyDefault
	} else {
		err := checkKeyIllegal(metadata[strategyKey])
		if err != nil {
			return err
		}
	}

	statesConfiguration[storeName] = &StoreConfiguration{keyPrefixStrategy: strategy}
	return nil
}

func GetModifiedStateKey(key, storeName, appID string) (string, error) {
	if err := checkKeyIllegal(key); err != nil {
		return "", err
	}
	stateConfiguration := getStateConfiguration(storeName)
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
	splits := strings.Split(modifiedStateKey, daprSeparator)
	if len(splits) <= 1 {
		return modifiedStateKey
	}
	return splits[1]
}

func getStateConfiguration(storeName string) *StoreConfiguration {
	c := statesConfiguration[storeName]
	if c == nil {
		c = &StoreConfiguration{keyPrefixStrategy: strategyDefault}
		statesConfiguration[storeName] = c
	}

	return c
}

func checkKeyIllegal(key string) error {
	if strings.Contains(key, daprSeparator) {
		return errors.Errorf("input key/keyPrefix '%s' can't contain '%s'", key, daprSeparator)
	}
	return nil
}

func StateConsistencyToString(c runtimev1pb.StateOptions_StateConsistency) string {
	switch c {
	case runtimev1pb.StateOptions_CONSISTENCY_EVENTUAL:
		return "eventual"
	case runtimev1pb.StateOptions_CONSISTENCY_STRONG:
		return "strong"
	}

	return ""
}

func StateConcurrencyToString(c runtimev1pb.StateOptions_StateConcurrency) string {
	switch c {
	case runtimev1pb.StateOptions_CONCURRENCY_FIRST_WRITE:
		return "first-write"
	case runtimev1pb.StateOptions_CONCURRENCY_LAST_WRITE:
		return "last-write"
	}

	return ""
}
