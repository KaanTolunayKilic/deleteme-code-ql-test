package mongo

import (
	"context"
	"fmt"
	"sync"

	"git.haw-hamburg.de/aci822/gabi/tweet-extractor/internal/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

type MongoDatabase struct {
	db      *mongo.Database
	Session *mongo.Client
	logger  *zap.Logger
}

func NewDatastore(config config.Configuration, logger *zap.Logger) *MongoDatabase {

	var mongoDataStore *MongoDatabase
	db, session := connect(config, logger)
	if db != nil && session != nil {
		mongoDataStore = new(MongoDatabase)
		mongoDataStore.db = db
		mongoDataStore.logger = logger
		mongoDataStore.Session = session
		return mongoDataStore
	}

	logger.Fatal(fmt.Sprintf("Failed to connect to database: %v", config.Mongo.Database))

	return nil
}

const mongPW = "test"

func connect(config config.Configuration, logger *zap.Logger) (a *mongo.Database, b *mongo.Client) {
	var connectOnce sync.Once
	var db *mongo.Database
	var session *mongo.Client
	connectOnce.Do(func() {
		db, session = connectToMongo(config, logger)
	})
	return db, session
}

func connectToMongo(config config.Configuration, logger *zap.Logger) (a *mongo.Database, b *mongo.Client) {

	var err error
	session, err := mongo.NewClient(options.Client().ApplyURI(fmt.Sprintf("mongodb://root:%s @localhost:27017", mongPW)))
	if err != nil {
		logger.Fatal(err.Error())
	}
	session.Connect(context.TODO())
	if err != nil {
		logger.Fatal(err.Error())
	}

	var DB = session.Database(config.Mongo.Database)
	logger.Info(fmt.Sprintf("Successfully connected to mongo database: %v", config.Mongo.Database))

	return DB, session
}
