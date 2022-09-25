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
package mock

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"mosn.io/layotto/components/pkg/utils"
)

type MockMongoFactory struct{}

// MockMongoClient is a mock of MongoClient interface
type MockMongoClient struct{}

// MockMongoSession is a mock of MongoSession interface
type MockMongoSession struct {
	mongo.SessionContext
}

// MockMongoCollection is a mock of MongoCollection interface
type MockMongoCollection struct {
	// '_id' document
	Result           map[string]bson.M
	InsertManyResult *mongo.InsertManyResult
	InsertOneResult  *mongo.InsertOneResult
	SingleResult     *mongo.SingleResult
	DeleteResult     *mongo.DeleteResult
}

func NewMockMongoFactory() *MockMongoFactory {
	return &MockMongoFactory{}
}

func NewMockMongoClient() *MockMongoClient {
	return &MockMongoClient{}
}

func NewMockMongoCollection() *MockMongoCollection {
	return &MockMongoCollection{}
}

func NewMockMongoSession() *MockMongoSession {
	return &MockMongoSession{}
}

func (f *MockMongoFactory) NewMongoClient(m utils.MongoMetadata) (utils.MongoClient, error) {
	return &MockMongoClient{}, nil
}

func (f *MockMongoFactory) NewMongoCollection(m *mongo.Database, collectionName string, opts *options.CollectionOptions) utils.MongoCollection {
	return &MockMongoCollection{}
}

func (f *MockMongoFactory) NewSingleResult(sr *mongo.SingleResult) utils.MongoSingleResult {
	return nil
}

func (mc *MockMongoCollection) FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult {
	result := mongo.SingleResult{}
	return &result
}

func (mc *MockMongoCollection) InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
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

func (mc *MockMongoCollection) DeleteOne(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	res := &mongo.DeleteResult{}
	doc := filter.(bson.M)
	value := doc["_id"].(string)
	if v, ok := mc.Result[value]; ok {
		if v["LockOwner"] == doc["LockOwner"] {
			delete(mc.Result, value)
			res.DeletedCount = 1
		}
	}
	return res, nil
}

func (mc *MockMongoCollection) Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error) {
	cursor := &mongo.Cursor{}
	return cursor, nil
}

func (mc *MockMongoCollection) Indexes() mongo.IndexView {
	return mongo.IndexView{}
}

func (mc *MockMongoCollection) UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	return nil, nil
}

func (mc *MockMongoCollection) FindOneAndUpdate(ctx context.Context, filter interface{}, update interface{}, opts ...*options.FindOneAndUpdateOptions) *mongo.SingleResult {
	return nil
}

func (c *MockMongoClient) StartSession(opts ...*options.SessionOptions) (mongo.Session, error) {
	return &MockMongoSession{}, nil
}

func (c *MockMongoClient) Ping(ctx context.Context, rp *readpref.ReadPref) error {
	return nil
}

func (c *MockMongoClient) Database(name string, opts ...*options.DatabaseOptions) *mongo.Database {
	return nil
}

func (c *MockMongoClient) Disconnect(ctx context.Context) error {
	return nil
}

func (s *MockMongoSession) AbortTransaction(context.Context) error {
	return nil
}

func (s *MockMongoSession) CommitTransaction(context.Context) error {
	return nil
}

func (s *MockMongoSession) WithTransaction(ctx context.Context, fn func(sessCtx mongo.SessionContext) (interface{}, error),
	opts ...*options.TransactionOptions) (interface{}, error) {
	res, err := fn(s)
	return res, err
}

func (s *MockMongoSession) EndSession(context.Context) {

}
