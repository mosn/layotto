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
	"reflect"
	"unsafe"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"mosn.io/layotto/components/pkg/utils"
)

var Result = make(map[string]bson.M)
var id string

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
	InsertOneResult *mongo.InsertOneResult
	Result          map[string]bson.M
}

type MockMongoSequencerSingleResult struct{}

func NewMockMongoSequencerFactory() *MockMongoSequencerFactory {
	return &MockMongoSequencerFactory{}
}

func NewMockMongoSequencerSession() *MockMongoSequencerSession {
	return &MockMongoSequencerSession{}
}

func (f *MockMongoSequencerFactory) NewSingleResult(sr *mongo.SingleResult) utils.MongoSingleResult {
	return &MockMongoSequencerSingleResult{}
}

func (f *MockMongoSequencerFactory) NewMongoClient(m utils.MongoMetadata) (utils.MongoClient, error) {
	return &MockMongoSequencerClient{}, nil
}

func (f *MockMongoSequencerFactory) NewMongoCollection(m *mongo.Database, collectionName string, opts *options.CollectionOptions) utils.MongoCollection {
	return &MockMongoSequencerCollection{}
}

func (mc *MockMongoSequencerCollection) FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult {
	result := new(mongo.SingleResult)
	value := reflect.ValueOf(result)
	doc := filter.(bson.M)
	id = doc["_id"].(string)
	if value.Kind() == reflect.Ptr {
		elem := value.Elem()
		err := elem.FieldByName("err")
		*(*error)(unsafe.Pointer(err.Addr().Pointer())) = nil

		cur := elem.FieldByName("cur")
		*(**mongo.Cursor)(unsafe.Pointer(cur.Addr().Pointer())) = &mongo.Cursor{}

		rdr := elem.FieldByName("rdr")
		*(*bson.Raw)(unsafe.Pointer(rdr.Addr().Pointer())) = bson.Raw{}

		reg := elem.FieldByName("reg")
		*(**bsoncodec.Registry)(unsafe.Pointer(reg.Addr().Pointer())) = &bsoncodec.Registry{}
	}
	return result
}

func (mc *MockMongoSequencerCollection) InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	doc := document.(bson.M)
	value := doc["_id"].(string)
	if _, ok := Result[value]; ok {
		return nil, nil
	} else {
		// insert cache
		Result[value] = doc
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
	if res, ok := Result[value]; ok {
		Result[value] = bson.M{"_id": res["_id"], "sequencer_value": res["sequencer_value"].(int) + 1}
		return nil, nil
	}
	return nil, nil
}

func (mc *MockMongoSequencerCollection) FindOneAndUpdate(ctx context.Context, filter interface{},
	update interface{}, opts ...*options.FindOneAndUpdateOptions) *mongo.SingleResult {
	doc := filter.(bson.M)
	id = doc["_id"].(string)
	upDoc := update.(bson.M)
	value := doc["_id"].(string)
	up := upDoc["$inc"].(bson.M)
	num := up["sequencer_value"].(int)
	if res, ok := Result[value]; ok {
		Result[value] = bson.M{"_id": res["_id"], "sequencer_value": res["sequencer_value"].(int) + num}
		return nil
	} else {
		bm := bson.M{"_id": doc["_id"], "sequencer_value": num}
		mc.InsertOne(ctx, bm)
	}
	return nil
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
	b := Result[id]
	id := b["_id"].(string)
	value := b["sequencer_value"].(int)
	ref := reflect.ValueOf(v)
	if ref.Kind() == reflect.Ptr {
		elem := ref.Elem()
		Id := elem.FieldByName("Id")
		*(*string)(unsafe.Pointer(Id.Addr().Pointer())) = id

		v := elem.FieldByName("Sequencer_value")
		*(*int)(unsafe.Pointer(v.Addr().Pointer())) = value
	}
	return nil
}

func (d *MockMongoSequencerSingleResult) Err() error {
	if Result[id] == nil {
		return mongo.ErrNoDocuments
	}
	return nil
}

func (d *MockMongoSequencerSingleResult) DecodeBytes() (bson.Raw, error) {
	return nil, nil
}
