//
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
package mongo

import (
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"mosn.io/layotto/components/pkg/utils"
	"mosn.io/layotto/components/sequencer"
	"mosn.io/pkg/log"
	"testing"
)

const key = "resource_xxx"

func TestMongoSequencer_Init(t *testing.T) {
	var mongoUrl = "localhost:xxxxx"
	comp := NewMongoSequencer(log.DefaultLogger)

	cfg := sequencer.Configuration{
		BiggerThan: nil,
		Properties: make(map[string]string),
	}

	cfg.Properties["mongoHost"] = mongoUrl
	err := comp.Init(cfg)

	assert.Error(t, err)
}

func TestMongoSequencer_GetNextId(t *testing.T) {
	var mongoUrl = "localhost:xxxxx"

	comp := NewMongoSequencer(log.DefaultLogger)

	cfg := sequencer.Configuration{
		BiggerThan: nil,
		Properties: make(map[string]string),
	}

	cfg.Properties["mongoHost"] = mongoUrl
	_ = comp.Init(cfg)

	// mock
	mockMongoClient := utils.MockMongoSequencerClient{}
	mockMongoSession := utils.NewMockMongoSequencerSession()
	mockMongoFactory := utils.NewMockMongoSequencerFactory()
	insertOneResult := &mongo.InsertOneResult{}
	mockMongoCollection := utils.MockMongoSequencerCollection{
		InsertOneResult: insertOneResult,
	}

	comp.session = mockMongoSession
	comp.collection = &mockMongoCollection
	comp.client = &mockMongoClient
	comp.factory = mockMongoFactory

	v1, err1 := comp.GetNextId(&sequencer.GetNextIdRequest{
		Key: key,
	})
	var expected1 int64 = 1
	assert.Equal(t, expected1, v1.NextId)
	assert.NoError(t, err1)

	v2, err2 := comp.GetNextId(&sequencer.GetNextIdRequest{
		Key: key,
	})
	var expected2 int64 = 2
	assert.Equal(t, expected2, v2.NextId)
	assert.NoError(t, err2)
}
