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
package utils

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MockMongoSequencerFactory struct{}

// MockMongoClient is a mock of MongoClient interface
type MockMongoSequencerClient struct{}

// MockMongoSession is a mock of MongoSession interface
type MockMongoSequencerSession struct {
	mongo.SessionContext
}

// MockMongoCollection is a mock of MongoCollection interface
type MockMongoSequencerCollection struct {
	// '_id' document
	Result          map[string]bson.M
	InsertOneResult *mongo.InsertOneResult
}

type MockMongoSequencerSingleResult struct {
	mongo.SingleResult
}

func NewMockMongoSequencerSession() *MockMongoSequencerSession {
	return &MockMongoSequencerSession{}
}

func (f *MockMongoSequencerFactory) NewMongoClient(m MongoMetadata) (MongoClient, error) {
	return &MockMongoSequencerClient{}, nil
}

func (f *MockMongoSequencerFactory) NewMongoCollection(m *mongo.Database, collectionName string, opts *options.CollectionOptions) MongoCollection {
	return &MockMongoSequencerCollection{}
}

func (mc *MockMongoSequencerCollection) FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult {
	res := mongo.SingleResult{}
	return &res
}

func (mc *MockMongoSequencerCollection) InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	doc := document.(bson.M)
	value := doc["_id"].(string)
	if _, ok := mc.Result[value]; ok {
		return nil, nil
	} else {
		// insert cache
		mc.Result[value] = doc
		mc.InsertOneResult.InsertedID = value
		return mc.InsertOneResult, nil
	}
}

func (mc *MockMongoSequencerCollection) DeleteOne(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	return nil, nil
}

func (mc *MockMongoSequencerCollection) Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error) {
	cursor := &mongo.Cursor{}
	return cursor, nil
}

func (mc *MockMongoSequencerCollection) Indexes() mongo.IndexView {
	return mongo.IndexView{}
}

func (mc *MockMongoSequencerCollection) UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	doc := filter.(bson.M)
	value := doc["_id"].(string)
	if res, ok := mc.Result[value]; ok {
		mc.Result[value] = bson.M{"_id": res["_id"], "sequencer_value": res["sequencer_value"].(int) + 1}
		return nil, nil
	}
	return nil, nil
}

func (c *MockMongoSequencerClient) StartSession(opts ...*options.SessionOptions) (mongo.Session, error) {
	return &MockMongoSequencerSession{}, nil
}

func (c *MockMongoSequencerClient) Ping(ctx context.Context, rp *readpref.ReadPref) error {
	return nil
}

func (c *MockMongoSequencerClient) Database(name string, opts ...*options.DatabaseOptions) *mongo.Database {
	return nil
}

func (c *MockMongoSequencerClient) Disconnect(ctx context.Context) error {
	return nil
}

func (s *MockMongoSequencerSession) AbortTransaction(context.Context) error {
	return nil
}

func (s *MockMongoSequencerSession) CommitTransaction(context.Context) error {
	return nil
}

func (s *MockMongoSequencerSession) WithTransaction(ctx context.Context, fn func(sessCtx mongo.SessionContext) (interface{}, error),
	opts ...*options.TransactionOptions) (interface{}, error) {
	res, err := fn(s)
	return res, err
}

func (s *MockMongoSequencerSession) EndSession(context.Context) {}

func (d *MockMongoSequencerSingleResult) Decode(v interface{}) error {
	return nil
}

func (d *MockMongoSequencerSingleResult) Err() error {
	return nil
}

func (d *MockMongoSequencerSingleResult) DecodeBytes() (bson.Raw, error) {
	return nil, nil
}
