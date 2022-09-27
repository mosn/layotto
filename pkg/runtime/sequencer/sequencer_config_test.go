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
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

const key = "lock-key-1234567"

func TestMain(m *testing.M) {
	SaveSeqConfiguration("store1", map[string]string{strategyKey: strategyNone})
	SaveSeqConfiguration("store2", map[string]string{strategyKey: strategyAppid})
	SaveSeqConfiguration("store3", map[string]string{strategyKey: strategyDefault})
	SaveSeqConfiguration("store4", map[string]string{strategyKey: strategyStoreName})
	SaveSeqConfiguration("store5", map[string]string{strategyKey: "other-fixed-prefix"})
	// if strategyKey not set
	SaveSeqConfiguration("store6", map[string]string{})
	os.Exit(m.Run())
}

func TestNonePrefix(t *testing.T) {
	modifiedLockKey, _ := GetModifiedSeqKey(key, "store1", "appid1")
	require.Equal(t, "sequencer|||"+key, modifiedLockKey)
}

func TestAppidPrefix(t *testing.T) {
	modifiedLockKey, _ := GetModifiedSeqKey(key, "store2", "appid1")
	require.Equal(t, "sequencer|||appid1||lock-key-1234567", modifiedLockKey)
}

func TestAppidPrefix_WithEnptyAppid(t *testing.T) {
	modifiedLockKey, _ := GetModifiedSeqKey(key, "store2", "")
	require.Equal(t, "sequencer|||lock-key-1234567", modifiedLockKey)
}

func TestDefaultPrefix(t *testing.T) {
	modifiedLockKey, _ := GetModifiedSeqKey(key, "store3", "appid1")
	require.Equal(t, "sequencer|||appid1||lock-key-1234567", modifiedLockKey)
}

func TestStoreNamePrefix(t *testing.T) {
	modifiedLockKey, _ := GetModifiedSeqKey(key, "store4", "appid1")
	require.Equal(t, "sequencer|||store4||lock-key-1234567", modifiedLockKey)
}

func TestOtherFixedPrefix(t *testing.T) {
	modifiedLockKey, _ := GetModifiedSeqKey(key, "store5", "appid1")
	require.Equal(t, "sequencer|||other-fixed-prefix||lock-key-1234567", modifiedLockKey)
}

func TestLegacyPrefix(t *testing.T) {
	modifiedLockKey, _ := GetModifiedSeqKey(key, "store6", "appid1")
	require.Equal(t, "sequencer|||appid1||lock-key-1234567", modifiedLockKey)
}

func TestPrefix_StoreNotInitial(t *testing.T) {
	// no config for store999
	modifiedLockKey, _ := GetModifiedSeqKey(key, "store999", "appid99")
	require.Equal(t, "sequencer|||appid99||lock-key-1234567", modifiedLockKey)
}
