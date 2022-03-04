/*
 * Copyright 2021 Layotto Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package state

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	SaveStateConfiguration("store1", map[string]string{strategyKey: strategyNone})
	SaveStateConfiguration("store2", map[string]string{strategyKey: strategyAppid})
	SaveStateConfiguration("store3", map[string]string{strategyKey: strategyDefault})
	SaveStateConfiguration("store4", map[string]string{strategyKey: strategyStoreName})
	SaveStateConfiguration("store5", map[string]string{strategyKey: "other-fixed-prefix"})
	// if strategyKey not set
	SaveStateConfiguration("store6", map[string]string{})
	os.Exit(m.Run())
}

func TestSaveStateConfiguration(t *testing.T) {
	testIllegalKeys := []struct {
		storename string
		prefix    string
	}{
		{
			storename: "statestore01",
			prefix:    "a||b",
		},
	}
	for _, item := range testIllegalKeys {
		err := SaveStateConfiguration(item.storename, map[string]string{
			strategyKey: item.prefix,
		})
		require.NotNil(t, err)
	}
}

func TestGetModifiedStateKey(t *testing.T) {
	// use custom prefix key
	testIllegalKeys := []struct {
		storename string
		prefix    string
		key       string
	}{
		{
			storename: "statestore01",
			prefix:    "a",
			key:       "c||d",
		},
	}
	for _, item := range testIllegalKeys {
		err := SaveStateConfiguration(item.storename, map[string]string{
			strategyKey: item.prefix,
		})
		require.Nil(t, err)
		_, err = GetModifiedStateKey(item.key, item.storename, "")
		require.NotNil(t, err)
	}
}

func TestNonePrefix(t *testing.T) {
	var key = "state-key-1234567"

	modifiedStateKey, _ := GetModifiedStateKey(key, "store1", "appid1")
	require.Equal(t, key, modifiedStateKey)

	originalStateKey := GetOriginalStateKey(modifiedStateKey)
	require.Equal(t, key, originalStateKey)
}

func TestAppidPrefix(t *testing.T) {
	var key = "state-key-1234567"

	modifiedStateKey, _ := GetModifiedStateKey(key, "store2", "appid1")
	require.Equal(t, "appid1||state-key-1234567", modifiedStateKey)

	originalStateKey := GetOriginalStateKey(modifiedStateKey)
	require.Equal(t, key, originalStateKey)
}

func TestAppidPrefix_WithEnptyAppid(t *testing.T) {
	var key = "state-key-1234567"

	modifiedStateKey, _ := GetModifiedStateKey(key, "store2", "")
	require.Equal(t, "state-key-1234567", modifiedStateKey)

	originalStateKey := GetOriginalStateKey(modifiedStateKey)
	require.Equal(t, key, originalStateKey)
}

func TestDefaultPrefix(t *testing.T) {
	var key = "state-key-1234567"

	modifiedStateKey, _ := GetModifiedStateKey(key, "store3", "appid1")
	require.Equal(t, "state-key-1234567", modifiedStateKey)

	originalStateKey := GetOriginalStateKey(modifiedStateKey)
	require.Equal(t, key, originalStateKey)
}

func TestStoreNamePrefix(t *testing.T) {
	var key = "state-key-1234567"

	modifiedStateKey, _ := GetModifiedStateKey(key, "store4", "appid1")
	require.Equal(t, "store4||state-key-1234567", modifiedStateKey)

	originalStateKey := GetOriginalStateKey(modifiedStateKey)
	require.Equal(t, key, originalStateKey)
}

func TestOtherFixedPrefix(t *testing.T) {
	var key = "state-key-1234567"

	modifiedStateKey, _ := GetModifiedStateKey(key, "store5", "appid1")
	require.Equal(t, "other-fixed-prefix||state-key-1234567", modifiedStateKey)

	originalStateKey := GetOriginalStateKey(modifiedStateKey)
	require.Equal(t, key, originalStateKey)
}

func TestLegacyPrefix(t *testing.T) {
	var key = "state-key-1234567"

	modifiedStateKey, _ := GetModifiedStateKey(key, "store6", "appid1")
	require.Equal(t, "state-key-1234567", modifiedStateKey)

	originalStateKey := GetOriginalStateKey(modifiedStateKey)
	require.Equal(t, key, originalStateKey)
}

func TestPrefix_StoreNotInitial(t *testing.T) {
	var key = "state-key-1234567"

	// no config for store999
	modifiedStateKey, _ := GetModifiedStateKey(key, "store999", "appid99")
	require.Equal(t, "state-key-1234567", modifiedStateKey)

	originalStateKey := GetOriginalStateKey(modifiedStateKey)
	require.Equal(t, key, originalStateKey)
}
