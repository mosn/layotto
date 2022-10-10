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
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"mosn.io/pkg/log"

	"mosn.io/layotto/components/pkg/mock"
	"mosn.io/layotto/components/sequencer"
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
	var mongoUrl = "localhost:27017"

	comp := NewMongoSequencer(log.DefaultLogger)

	cfg := sequencer.Configuration{
		BiggerThan: nil,
		Properties: make(map[string]string),
	}

	cfg.Properties["mongoHost"] = mongoUrl
	_ = comp.Init(cfg)

	// mock
	mockMongoClient := mock.MockMongoSequencerClient{}
	mockMongoSession := mock.NewMockMongoSequencerSession()
	mockMongoFactory := mock.NewMockMongoSequencerFactory()
	insertOneResult := &mongo.InsertOneResult{}
	mockMongoCollection := mock.MockMongoSequencerCollection{
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
	delete(mock.Result, key)
	var expected2 int64 = 2
	assert.Equal(t, expected2, v2.NextId)
	assert.NoError(t, err2)
}

func TestMongoSequencer_Close(t *testing.T) {
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
	mockMongoClient := mock.MockMongoSequencerClient{}
	mockMongoSession := mock.NewMockMongoSequencerSession()
	mockMongoFactory := mock.NewMockMongoSequencerFactory()
	insertOneResult := &mongo.InsertOneResult{}
	mockMongoCollection := mock.MockMongoSequencerCollection{
		InsertOneResult: insertOneResult,
		Result:          result,
	}

	comp.session = mockMongoSession
	comp.collection = &mockMongoCollection
	comp.client = &mockMongoClient
	comp.factory = mockMongoFactory

	err := comp.Close()
	assert.NoError(t, err)
}

func TestMongoSequencer_GetSegment(t *testing.T) {
	var mongoUrl = "localhost:xxxxx"

	comp := NewMongoSequencer(log.DefaultLogger)

	cfg := sequencer.Configuration{
		BiggerThan: nil,
		Properties: make(map[string]string),
	}

	cfg.Properties["mongoHost"] = mongoUrl
	_ = comp.Init(cfg)

	// mock
	mockMongoClient := mock.MockMongoSequencerClient{}
	mockMongoSession := mock.NewMockMongoSequencerSession()
	mockMongoFactory := mock.NewMockMongoSequencerFactory()
	insertOneResult := &mongo.InsertOneResult{}
	mockMongoCollection := mock.MockMongoSequencerCollection{
		InsertOneResult: insertOneResult,
	}

	comp.session = mockMongoSession
	comp.collection = &mockMongoCollection
	comp.client = &mockMongoClient
	comp.factory = mockMongoFactory

	support, res, err := comp.GetSegment(&sequencer.GetSegmentRequest{
		Size: 20,
		Key:  key,
	})
	var expected1 int64 = 1
	var expected2 int64 = 20
	assert.Equal(t, support, true)
	assert.Equal(t, expected1, res.From)
	assert.Equal(t, expected2, res.To)
	assert.NoError(t, err)
}
