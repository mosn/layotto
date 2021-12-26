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
	"go.mongodb.org/mongo-driver/bson"
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
	result := make(map[string]bson.M)
	mockMongoClient := utils.MockMongoSequencerClient{}
	mockMongoSession := utils.NewMockMongoSequencerSession()
	insertOneResult := &mongo.InsertOneResult{}
	mockMongoCollection := utils.MockMongoSequencerCollection{
		Result:          result,
		InsertOneResult: insertOneResult,
	}

	comp.session = mockMongoSession
	comp.collection = &mockMongoCollection
	comp.client = &mockMongoClient

	v, err := comp.GetNextId(&sequencer.GetNextIdRequest{
		Key: key,
	})
	var expected int64 = 0
	assert.Equal(t, expected, v.NextId)
	assert.NoError(t, err)
}
