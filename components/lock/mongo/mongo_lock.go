package mongo

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"mosn.io/layotto/components/lock"
	"mosn.io/layotto/components/pkg/utils"
	"mosn.io/pkg/log"
	"time"
)

// mongo lock store
type MongoLock struct {
	client     *mongo.Client
	collection *mongo.Collection
	database   *mongo.Database
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
	// 1.parse config
	m, err := utils.ParseMongoMetadata(metadata.Properties)
	if err != nil {
		return err
	}
	e.metadata = m

	// 2. construct client
	if e.client, err = utils.NewMongoClient(m); err != nil {
		return err
	}

	e.ctx, e.cancel = context.WithCancel(context.Background())

	if err := e.client.Ping(e.ctx, nil); err != nil {
		return err
	}

	wc, err := utils.GetWriteConcernObject(e.metadata.WriteConcern)
	if err != nil {
		return err
	}

	rc, err := utils.GetReadConcrenObject(e.metadata.ReadConcern)
	if err != nil {
		return err
	}

	opts := options.Collection().SetWriteConcern(wc).SetReadConcern(rc)
	database := e.client.Database(e.metadata.DatabaseName)
	collection := database.Collection(e.metadata.CollectionName, opts)
	e.database = database
	e.collection = collection

	// create exprie time index
	indexModel := mongo.IndexModel{
		Keys:    bsonx.Doc{{"Expire", bsonx.Int64(1)}},
		Options: options.Index().SetExpireAfterSeconds(0),
	}
	e.collection.Indexes().CreateOne(e.ctx, indexModel)

	return err
}

// Features is to get MongoLock's features
func (e *MongoLock) Features() []lock.Feature {
	return e.features
}

func (e *MongoLock) TryLock(req *lock.TryLockRequest) (*lock.TryLockResponse, error) {
	session, err := e.client.StartSession()
	txnOpts := options.Transaction().SetReadConcern(readconcern.Snapshot()).
		SetWriteConcern(writeconcern.New(writeconcern.WMajority()))

	defer session.EndSession(e.ctx)

	// check session
	if err != nil {
		return &lock.TryLockResponse{}, fmt.Errorf("[mongoLock]: Create session return error: %s ResourceId: %s", err, req.ResourceId)
	}

	status, err := session.WithTransaction(e.ctx, func(sessionContext mongo.SessionContext) (interface{}, error) {
		var status int
		var err error
		var singleResult *mongo.SingleResult
		var insertOneResult *mongo.InsertOneResult

		// set exprie date
		expireTime := time.Now().Add(time.Duration(req.Expire) * time.Second)

		// find mongo lock
		singleResult = e.collection.FindOne(e.ctx, bson.M{"_id": req.ResourceId})

		// insert mongo lock
		if singleResult.Err() == mongo.ErrNoDocuments {
			insertOneResult, err = e.collection.InsertOne(e.ctx, bson.M{"_id": req.ResourceId, "LockOwner": req.LockOwner, "Expire": expireTime})
		}

		if err != nil {
			sessionContext.AbortTransaction(sessionContext)
			return status, err
		}

		// commit and set status
		if insertOneResult != nil && insertOneResult.InsertedID == req.ResourceId {
			if err = sessionContext.CommitTransaction(sessionContext); err == nil {
				status = 1
			}
		}

		return status, err
	}, txnOpts)

	// check lock
	if err != nil {
		return &lock.TryLockResponse{}, fmt.Errorf("[mongoLock]: Create new lock return error: %s ResourceId: %s", err, req.ResourceId)
	}

	if status == 1 {
		return &lock.TryLockResponse{
			Success: true,
		}, nil
	} else {
		return &lock.TryLockResponse{
			Success: false,
		}, nil
	}
}

func (e *MongoLock) Unlock(req *lock.UnlockRequest) (*lock.UnlockResponse, error) {
	session, err := e.client.StartSession()
	txnOpts := options.Transaction().SetReadConcern(readconcern.Snapshot()).
		SetWriteConcern(writeconcern.New(writeconcern.WMajority()))

	defer session.EndSession(e.ctx)

	if err != nil {
		return newInternalErrorUnlockResponse(), fmt.Errorf("[mongoLock]: Create Session return error: %s ResourceId: %s", err, req.ResourceId)
	}

	status, err := session.WithTransaction(e.ctx, func(sessionContext mongo.SessionContext) (interface{}, error) {
		var status int

		result, err := e.collection.DeleteOne(e.ctx, bson.M{"_id": req.ResourceId, "LockOwner": req.LockOwner})

		if result.DeletedCount == 1 && err == nil {
			status = 0
		} else if result.DeletedCount == 0 && err == nil {
			if cursor, err := e.collection.Find(e.ctx, bson.M{"_id": req.ResourceId}); cursor != nil && cursor.RemainingBatchLength() != 0 && err == nil {
				status = 2
			} else if cursor != nil && cursor.RemainingBatchLength() == 0 && err == nil {
				status = 1
			}
		}

		if err != nil {
			sessionContext.AbortTransaction(sessionContext)
			return nil, err
		}
		err = sessionContext.CommitTransaction(sessionContext)

		return status, err
	}, txnOpts)

	if err != nil {
		return newInternalErrorUnlockResponse(), fmt.Errorf("[mongoLock]: Unlock returned error: %s ResourceId: %s", err, req.ResourceId)
	}

	resp := lock.INTERNAL_ERROR

	if status == 0 {
		resp = lock.SUCCESS
	} else if status == 1 {
		resp = lock.LOCK_UNEXIST
	} else if status == 2 {
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
