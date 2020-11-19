package mongodbc

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	// Environment Variables
	mongoDBServerHost     = "MONGODB_SERVER_HOST"
	mongoDBServerPort     = "MONGODB_SERVER_PORT"
	mongoDBContextTimeout = "MONGODB_CONTEXT_TIMEOUT"
	// Errors
	mongoDBEmptyHostAddrErr     = "DB Host address is empty"
	mongoDBEmptyPortErr         = "DB Port is empty"
	mongoDBCtxTimeoutAtoiErr    = "unable to read db context timeout"
	mongoDBCtxTimeoutValErr     = "value is less than or equal to 0"
	mongoDBDisconnectErrWrapper = "MongoClient Disconnect: %w"
)

// MongoConfig contains mongodb server
// configuration needed by the client
// in order to connect to it
type MongoConfig struct {
	Host   string
	Port   string
	URI    string
	Ctx    context.Context
	Cancel context.CancelFunc
}

// ReadEnv reads the necessary mongodb environment variables
func (c *MongoConfig) ReadEnv() error {
	c.Host = os.Getenv(mongoDBServerHost)
	if c.Host == "" {
		return errors.New(mongoDBEmptyHostAddrErr)
	}

	c.Port = os.Getenv(mongoDBServerPort)
	if c.Port == "" {
		return errors.New(mongoDBEmptyPortErr)
	}

	c.URI = "mongodb://" + c.Host + ":" + c.Port

	ctxTime, err := strconv.Atoi(os.Getenv(mongoDBContextTimeout))
	if err != nil {
		return errors.New(mongoDBCtxTimeoutAtoiErr)
	}

	if ctxTime <= 0 {
		return errors.New(mongoDBCtxTimeoutValErr)
	}

	duration := time.Duration(ctxTime) * time.Second
	c.Ctx, c.Cancel = context.WithTimeout(context.Background(), duration)

	return nil
}

// MongoClient is a mongodb client implementation
type MongoClient struct {
	Config MongoConfig
	Client *mongo.Client
}

// Connect to mongodb server
func (m *MongoClient) Connect() (err error) {
	defer m.Config.Cancel()
	log.Println("connecting to mongodb")
	m.Client, err = mongo.Connect(m.Config.Ctx, options.Client().ApplyURI(m.Config.URI))
	log.Println("connected to mongodb")
	return err
}

// Disconnect client from mongodb server
func (m *MongoClient) Disconnect() error {
	defer m.Config.Cancel()
	log.Println("disconnecting to mongodb")
	return fmt.Errorf(mongoDBDisconnectErrWrapper, m.Client.Disconnect(m.Config.Ctx))
}

// Ping the mongodb server
func (m *MongoClient) Ping() error {
	return m.Client.Ping(m.Config.Ctx, readpref.Primary())
}
