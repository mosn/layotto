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
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"mosn.io/pkg/log"

	"mosn.io/layotto/components/lock"
	"mosn.io/layotto/components/pkg/utils"
)

const (
	TRY_LOCK_SUCCESS        = 1
	TRY_LOCK_FAIL           = 2
	UNLOCK_SUCCESS          = 3
	UNLOCK_UNEXIST          = 4
	UNLOCK_BELONG_TO_OTHERS = 5
	UNLOCK_FAIL             = 6
)

// mongo lock store
type MongoLock struct {
	factory utils.MongoFactory

	client     utils.MongoClient
	session    utils.MongoSession
	collection utils.MongoCollection
	metadata   utils.MongoMetadata

	features []lock.Feature
	logger   log.ErrorLogger

	ctx    context.Context
	cancel context.CancelFunc
}

// NewMongoLock returns a new mongo lock
func NewMongoLock(logger log.ErrorLogger) *MongoLock {
	s := &MongoLock{
		features: make([]lock.Feature, 0),
		logger:   logger,
	}
	return s
}

func (e *MongoLock) Init(metadata lock.Metadata) error {
	var client utils.MongoClient
	// 1.parse config
	m, err := utils.ParseMongoMetadata(metadata.Properties)
	if err != nil {
		return err
	}
	e.metadata = m

	e.factory = &utils.MongoFactoryImpl{}

	// 2. construct client
	if client, err = e.factory.NewMongoClient(m); err != nil {
		return err
	}

	e.ctx, e.cancel = context.WithCancel(context.Background())

	if err := client.Ping(e.ctx, nil); err != nil {
		return err
	}
	// Connections Collection
	e.collection, err = utils.SetCollection(client, e.factory, e.metadata)
	if err != nil {
		return err
	}

	// create exprie time index
	indexModel := mongo.IndexModel{
		Keys:    bsonx.Doc{{Key: "Expire", Value: bsonx.Int64(1)}},
		Options: options.Index().SetExpireAfterSeconds(0),
	}
	e.collection.Indexes().CreateOne(e.ctx, indexModel)

	e.client = client

	return err
}

// Features is to get MongoLock's features
func (e *MongoLock) Features() []lock.Feature {
	return e.features
}

// LockKeepAlive try to renewal lease
func (e *MongoLock) LockKeepAlive(ctx context.Context, request *lock.LockKeepAliveRequest) (*lock.LockKeepAliveResponse, error) {
	//TODO: implemnt function
	return nil, nil
}

func (e *MongoLock) TryLock(req *lock.TryLockRequest) (*lock.TryLockResponse, error) {
	var err error
	// create mongo session
	e.session, err = e.client.StartSession()
	txnOpts := options.Transaction().SetReadConcern(readconcern.Snapshot()).
		SetWriteConcern(writeconcern.New(writeconcern.WMajority()))

	// check session
	if err != nil {
		return &lock.TryLockResponse{
			Success: false,
		}, fmt.Errorf("[mongoLock]: Create session return error: %s ResourceId: %s", err, req.ResourceId)
	}

	// close mongo session
	defer e.session.EndSession(e.ctx)

	// start transaction
	status, err := e.session.WithTransaction(e.ctx, func(sessionContext mongo.SessionContext) (interface{}, error) {
		var err error
		var insertOneResult *mongo.InsertOneResult

		// set exprie date
		expireTime := time.Now().Add(time.Duration(req.Expire) * time.Second)

		// insert mongo lock
		insertOneResult, err = e.collection.InsertOne(e.ctx, bson.M{"_id": req.ResourceId, "LockOwner": req.LockOwner, "Expire": expireTime})

		if err != nil {
			_ = sessionContext.AbortTransaction(sessionContext)
			return TRY_LOCK_FAIL, err
		}

		// commit and set status
		if insertOneResult != nil && insertOneResult.InsertedID == req.ResourceId {
			if err = sessionContext.CommitTransaction(sessionContext); err == nil {
				return TRY_LOCK_SUCCESS, err
			}
		}

		return TRY_LOCK_FAIL, err
	}, txnOpts)

	// check lock
	if err != nil {
		return &lock.TryLockResponse{}, fmt.Errorf("[mongoLock]: Create new lock return error: %s ResourceId: %s", err, req.ResourceId)
	}

	if status == TRY_LOCK_SUCCESS {
		return &lock.TryLockResponse{
			Success: true,
		}, nil
	}
	return &lock.TryLockResponse{
		Success: false,
	}, nil
}

func (e *MongoLock) Unlock(req *lock.UnlockRequest) (*lock.UnlockResponse, error) {
	var err error
	// create mongo session
	e.session, err = e.client.StartSession()
	txnOpts := options.Transaction().SetReadConcern(readconcern.Snapshot()).
		SetWriteConcern(writeconcern.New(writeconcern.WMajority()))

	// check session
	if err != nil {
		return newInternalErrorUnlockResponse(), fmt.Errorf("[mongoLock]: Create Session return error: %s ResourceId: %s", err, req.ResourceId)
	}

	// close mongo session
	defer e.session.EndSession(e.ctx)

	// start transaction
	status, err := e.session.WithTransaction(e.ctx, func(sessionContext mongo.SessionContext) (interface{}, error) {
		var status int

		// delete lock
		result, err := e.collection.DeleteOne(e.ctx, bson.M{"_id": req.ResourceId, "LockOwner": req.LockOwner})

		// check delete result
		if result.DeletedCount == 1 && err == nil {
			status = UNLOCK_SUCCESS
		} else if result.DeletedCount == 0 && err == nil {
			if cursor, err := e.collection.Find(e.ctx, bson.M{"_id": req.ResourceId}); cursor != nil && cursor.RemainingBatchLength() != 0 && err == nil {
				status = UNLOCK_BELONG_TO_OTHERS
			} else if cursor != nil && cursor.RemainingBatchLength() == 0 && err == nil {
				status = UNLOCK_UNEXIST
			}
		}

		if err != nil {
			_ = sessionContext.AbortTransaction(sessionContext)
			return UNLOCK_FAIL, err
		}

		// commit and set status
		if err = sessionContext.CommitTransaction(sessionContext); err == nil {
			return status, err
		}

		return UNLOCK_FAIL, err
	}, txnOpts)

	resp := lock.INTERNAL_ERROR

	if err != nil {
		return newInternalErrorUnlockResponse(), fmt.Errorf("[mongoLock]: Unlock returned error: %s ResourceId: %s", err, req.ResourceId)
	}

	if status == UNLOCK_SUCCESS {
		resp = lock.SUCCESS
	} else if status == UNLOCK_UNEXIST {
		resp = lock.LOCK_UNEXIST
	} else if status == UNLOCK_BELONG_TO_OTHERS {
		resp = lock.LOCK_BELONG_TO_OTHERS
	}
	return &lock.UnlockResponse{
		Status: resp,
	}, nil
}

func newInternalErrorUnlockResponse() *lock.UnlockResponse {
	return &lock.UnlockResponse{
		Status: lock.INTERNAL_ERROR,
	}
}

func (e *MongoLock) Close() error {
	e.cancel()

	return e.client.Disconnect(e.ctx)
}
