package utils

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"strconv"
	"time"
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

func NewMongoClient(m MongoMetadata) (*mongo.Client, error) {
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

	if val, ok := properties[mongoHost]; ok && val != "" {
		m.Host = val
	}

	if val, ok := properties[server]; ok && val != "" {
		m.Server = val
	}

	if len(m.Host) == 0 && len(m.Server) == 0 {
		return m, errors.New("must set 'host' or 'server' fields")
	}

	if len(m.Host) != 0 && len(m.Server) != 0 {
		return m, errors.New("'host' or 'server' fields are mutually exclusive")
	}

	if val, ok := properties[username]; ok && val != "" {
		m.Username = val
	}

	if val, ok := properties[mongoPassword]; ok && val != "" {
		m.Password = val
	}

	m.DatabaseName = defaultDatabase
	if val, ok := properties[databaseName]; ok && val != "" {
		m.DatabaseName = val
	}

	m.CollectionName = defaultCollectionName
	if val, ok := properties[collecttionName]; ok && val != "" {
		m.CollectionName = val
	}

	if val, ok := properties[writeConcern]; ok && val != "" {
		m.WriteConcern = val
	}

	if val, ok := properties[readConcern]; ok && val != "" {
		m.ReadConcern = val
	}

	if val, ok := properties[params]; ok && val != "" {
		m.Params = val
	}

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
