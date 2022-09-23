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

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"mosn.io/pkg/log"

	"mosn.io/layotto/components/pkg/utils"
	"mosn.io/layotto/components/sequencer"
)

type MongoSequencer struct {
	factory utils.MongoFactory

	client     utils.MongoClient
	session    utils.MongoSession
	collection utils.MongoCollection
	singResult utils.MongoSingleResult
	metadata   utils.MongoMetadata
	biggerThan map[string]int64

	logger log.ErrorLogger

	ctx    context.Context
	cancel context.CancelFunc
}

type SequencerDocument struct {
	Id              string `bson:"_id"`
	Sequencer_value int64  `bson:"sequencer_value"`
}

// MongoSequencer returns a new mongo sequencer
func NewMongoSequencer(logger log.ErrorLogger) *MongoSequencer {
	m := &MongoSequencer{
		logger: logger,
	}

	return m
}

func (e *MongoSequencer) Init(config sequencer.Configuration) error {
	var document SequencerDocument
	// 1.parse config
	m, err := utils.ParseMongoMetadata(config.Properties)
	if err != nil {
		return err
	}
	e.metadata = m
	e.biggerThan = config.BiggerThan

	e.factory = &utils.MongoFactoryImpl{}

	// 2. construct client
	e.ctx, e.cancel = context.WithCancel(context.Background())

	if e.client, err = e.factory.NewMongoClient(m); err != nil {
		return err
	}

	if err := e.client.Ping(e.ctx, nil); err != nil {
		return err
	}

	// Connections Collection
	e.collection, err = utils.SetCollection(e.client, e.factory, e.metadata)
	if err != nil {
		return err
	}

	if len(e.biggerThan) > 0 {
		for k, bt := range e.biggerThan {
			if bt <= 0 {
				continue
			}
			// find key of biggerThan
			cursor, err := e.collection.Find(e.ctx, bson.M{"_id": k})
			if err != nil {
				return err
			}
			if cursor != nil && cursor.RemainingBatchLength() > 0 {
				cursor.Decode(&document)
			}
			// check biggerThan's value
			if document.Sequencer_value < bt {
				return fmt.Errorf("mongo sequencer error: can not satisfy biggerThan guarantee.key: %s,current id:%v", k, document.Sequencer_value)
			}
		}
	}

	return err
}

func (e *MongoSequencer) GetNextId(req *sequencer.GetNextIdRequest) (*sequencer.GetNextIdResponse, error) {
	var err error
	var document SequencerDocument
	// create mongo session
	e.session, err = e.client.StartSession()
	txnOpts := options.Transaction().SetReadConcern(readconcern.Snapshot()).
		SetWriteConcern(writeconcern.New(writeconcern.WMajority()))

	// check session
	if err != nil {
		return nil, fmt.Errorf("[mongoSequencer]: Create Session return error: %s key: %s", err, req.Key)
	}

	// close mongo session
	defer e.session.EndSession(e.ctx)

	status, err := e.session.WithTransaction(e.ctx, func(sessionContext mongo.SessionContext) (interface{}, error) {
		var err error
		after := options.After
		upsert := true
		opt := options.FindOneAndUpdateOptions{
			ReturnDocument: &after,
			Upsert:         &upsert,
		}

		e.singResult = e.factory.NewSingleResult(e.collection.FindOneAndUpdate(e.ctx, bson.M{"_id": req.Key}, bson.M{"$inc": bson.M{"sequencer_value": 1}}, &opt))

		// rollback
		if e.singResult.Err() != nil {
			_ = sessionContext.AbortTransaction(sessionContext)
			return nil, err
		}

		// commit
		if err = sessionContext.CommitTransaction(sessionContext); err != nil {
			return nil, err
		}
		e.singResult.Decode(&document)
		return document, nil
	}, txnOpts)
	if err != nil || status == nil {
		return nil, err
	}

	return &sequencer.GetNextIdResponse{
		NextId: document.Sequencer_value,
	}, nil
}

func (e *MongoSequencer) GetSegment(req *sequencer.GetSegmentRequest) (support bool, result *sequencer.GetSegmentResponse, err error) {
	var document SequencerDocument

	// size=0 only check support
	if req.Size == 0 {
		return true, nil, nil
	}

	// create mongo session
	e.session, err = e.client.StartSession()
	txnOpts := options.Transaction().SetReadConcern(readconcern.Snapshot()).
		SetWriteConcern(writeconcern.New(writeconcern.WMajority()))

	// check session
	if err != nil {
		return true, nil, fmt.Errorf("[mongoSequencer]: Create Session return error: %s key: %s", err, req.Key)
	}

	// close mongo session
	defer e.session.EndSession(e.ctx)

	status, err := e.session.WithTransaction(e.ctx, func(sessionContext mongo.SessionContext) (interface{}, error) {
		var err error
		after := options.After
		upsert := true
		opt := options.FindOneAndUpdateOptions{
			ReturnDocument: &after,
			Upsert:         &upsert,
		}

		// find and upsert
		e.singResult = e.factory.NewSingleResult(e.collection.FindOneAndUpdate(e.ctx, bson.M{"_id": req.Key}, bson.M{"$inc": bson.M{"sequencer_value": req.Size}}, &opt))

		// rollback
		if e.singResult.Err() != nil {
			_ = sessionContext.AbortTransaction(sessionContext)
			return nil, err
		}

		// commit
		if err = sessionContext.CommitTransaction(sessionContext); err != nil {
			return nil, err
		}
		e.singResult.Decode(&document)
		return document, nil
	}, txnOpts)
	if err != nil || status == nil {
		return true, nil, err
	}

	return true, &sequencer.GetSegmentResponse{
		From: document.Sequencer_value - int64(req.Size) + 1,
		To:   document.Sequencer_value,
	}, nil
}

func (e *MongoSequencer) Close() error {
	e.cancel()

	return e.client.Disconnect(e.ctx)
}
