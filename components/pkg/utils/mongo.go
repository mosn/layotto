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
	"errors"
	"fmt"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

const (
	mongoHost        = "mongoHost"
	mongoPassword    = "mongoPassword"
	username         = "username"
	server           = "server"
	databaseName     = "databaseName"
	collecttionName  = "collectionName"
	writeConcern     = "writeConcern"
	readConcern      = "readConcern"
	operationTimeout = "operationTimeout"
	params           = "params"

	defaultDatabase       = "layottoStore"
	defaultCollectionName = "layottoCollection"
	defaultTimeout        = 5 * time.Second

	// mongodb://<username>:<password@<host>/<database><params>
	connectionURIFormatWithAuthentication = "mongodb://%s:%s@%s/%s%s"

	// mongodb://<host>/<database><params>
	connectionURIFormat = "mongodb://%s/%s%s"

	// mongodb+srv://<server>/<params>
	connectionURIFormatWithSrv = "mongodb+srv://%s/%s"
)

type MongoMetadata struct {
	Host             string
	Username         string
	Password         string
	DatabaseName     string
	CollectionName   string
	Server           string
	Params           string
	WriteConcern     string
	ReadConcern      string
	OperationTimeout time.Duration
}

// Item is Mongodb document wrapper.
type Item struct {
	Key   string      `bson:"_id"`
	Value interface{} `bson:"value"`
	Etag  string      `bson:"_etag"`
}

type MongoFactory interface {
	NewMongoClient(m MongoMetadata) (MongoClient, error)
	NewMongoCollection(m *mongo.Database, collectionName string, opts *options.CollectionOptions) MongoCollection
	NewSingleResult(sr *mongo.SingleResult) MongoSingleResult
}

type MongoClient interface {
	StartSession(opts ...*options.SessionOptions) (mongo.Session, error)
	Ping(ctx context.Context, rp *readpref.ReadPref) error
	Database(name string, opts ...*options.DatabaseOptions) *mongo.Database
	Disconnect(ctx context.Context) error
}

type MongoSession interface {
	AbortTransaction(context.Context) error
	CommitTransaction(context.Context) error
	WithTransaction(ctx context.Context, fn func(sessCtx mongo.SessionContext) (interface{}, error),
		opts ...*options.TransactionOptions) (interface{}, error)
	EndSession(context.Context)
}

type MongoCollection interface {
	FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult
	InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error)
	DeleteOne(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error)
	Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error)
	Indexes() mongo.IndexView
	UpdateOne(ctx context.Context, filter interface{}, update interface{},
		opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)
	FindOneAndUpdate(ctx context.Context, filter interface{},
		update interface{}, opts ...*options.FindOneAndUpdateOptions) *mongo.SingleResult
}

type MongoSingleResult interface {
	Decode(v interface{}) error
	Err() error
	DecodeBytes() (bson.Raw, error)
}

type MongoFactoryImpl struct{}

func (c *MongoFactoryImpl) NewSingleResult(sr *mongo.SingleResult) MongoSingleResult {
	return sr
}

func (c *MongoFactoryImpl) NewMongoCollection(m *mongo.Database, collectionName string, opts *options.CollectionOptions) MongoCollection {
	collection := m.Collection(collectionName, opts)
	return collection
}

func (c *MongoFactoryImpl) NewMongoClient(m MongoMetadata) (MongoClient, error) {
	uri := getMongoURI(m)

	// Set client options
	clientOptions := options.Client().ApplyURI(uri)

	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), m.OperationTimeout)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	return client, err
}

func ParseMongoMetadata(properties map[string]string) (MongoMetadata, error) {
	m := MongoMetadata{}

	m.Host = getString(properties, mongoHost)

	m.Server = getString(properties, server)

	if len(m.Host) == 0 && len(m.Server) == 0 {
		return m, errors.New("must set 'host' or 'server' fields")
	}

	if len(m.Host) != 0 && len(m.Server) != 0 {
		return m, errors.New("'host' or 'server' fields are mutually exclusive")
	}

	m.Username = getString(properties, username)
	m.Password = getString(properties, mongoPassword)

	m.DatabaseName = defaultDatabase
	if val, ok := properties[databaseName]; ok && val != "" {
		m.DatabaseName = val
	}

	m.CollectionName = defaultCollectionName
	if val, ok := properties[collecttionName]; ok && val != "" {
		m.CollectionName = val
	}

	m.WriteConcern = getString(properties, writeConcern)
	m.ReadConcern = getString(properties, readConcern)
	m.Params = getString(properties, params)

	var err error
	m.OperationTimeout = defaultTimeout
	if val, ok := properties[operationTimeout]; ok && val != "" {
		m.OperationTimeout, err = time.ParseDuration(val)
		if err != nil {
			return m, errors.New("incorrect operationTimeout field")
		}
	}
	return m, nil
}

func getString(properties map[string]string, key string) string {
	if val, ok := properties[key]; ok && val != "" {
		return val
	}
	return ""
}

func getMongoURI(m MongoMetadata) string {
	if len(m.Server) != 0 {
		return fmt.Sprintf(connectionURIFormatWithSrv, m.Server, m.Params)
	}

	if m.Username != "" && m.Password != "" {
		return fmt.Sprintf(connectionURIFormatWithAuthentication, m.Username, m.Password, m.Host, m.DatabaseName, m.Params)
	}

	return fmt.Sprintf(connectionURIFormat, m.Host, m.DatabaseName, m.Params)
}

func GetWriteConcernObject(cn string) (*writeconcern.WriteConcern, error) {
	var wc *writeconcern.WriteConcern
	if cn != "" {
		if cn == "majority" {
			wc = writeconcern.New(writeconcern.WMajority(), writeconcern.J(true), writeconcern.WTimeout(defaultTimeout))
		} else {
			w, err := strconv.Atoi(cn)
			wc = writeconcern.New(writeconcern.W(w), writeconcern.J(true), writeconcern.WTimeout(defaultTimeout))
			return wc, err
		}
	} else {
		wc = writeconcern.New(writeconcern.W(1), writeconcern.J(true), writeconcern.WTimeout(defaultTimeout))
	}
	return wc, nil
}

func GetReadConcrenObject(cn string) (*readconcern.ReadConcern, error) {
	switch cn {
	case "local":
		return readconcern.Local(), nil
	case "majority":
		return readconcern.Majority(), nil
	case "available":
		return readconcern.Available(), nil
	case "linearizable":
		return readconcern.Linearizable(), nil
	case "snapshot":
		return readconcern.Snapshot(), nil
	case "":
		return readconcern.Local(), nil
	}
	return nil, nil
}

func SetConcern(m MongoMetadata) (*options.CollectionOptions, error) {

	wc, err := GetWriteConcernObject(m.WriteConcern)
	if err != nil {
		return nil, err
	}

	rc, err := GetReadConcrenObject(m.ReadConcern)
	if err != nil {
		return nil, err
	}

	// set mongo options of collection
	opts := options.Collection().SetWriteConcern(wc).SetReadConcern(rc)

	return opts, err
}

func SetCollection(c MongoClient, f MongoFactory, m MongoMetadata) (MongoCollection, error) {
	opts, err := SetConcern(m)
	if err != nil {
		return nil, err
	}
	return f.NewMongoCollection(c.Database(m.DatabaseName), m.CollectionName, opts), err
}
