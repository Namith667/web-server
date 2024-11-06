package service

import (
	"context"
	"errors"
	"os"
	"web-server/logger"
	"web-server/model"

	"go.uber.org/zap"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//	type Database struct {
//		Client *mongo.Client
//	}
//
// Database is a wrapper for the MongoDB client
type Database struct {
	Client *mongo.Client
	logg   *logger.Logger
}

func InitDatabase(logg *logger.Logger) (*Database, error) {
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		logg.Error("Set the 'MONGODB_URI' environment variable.")
		return nil, errors.New("no MongoDB URI provided")
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		logg.Error("Failed to connect to MongoDB", zap.Error(err))
		return nil, err
	}

	logg.Info("Connected to MongoDB")
	return &Database{Client: client, logg: logg}, nil
}

///

func (db *Database) Disconnect() error {
	return db.Client.Disconnect(context.TODO())
}

func (db *Database) Create(collectionName string, user model.User) (*mongo.InsertOneResult, error) {
	coll := db.Client.Database("sample_db").Collection(collectionName)
	result, err := coll.InsertOne(context.TODO(), user)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (db *Database) Read(collectionName string) ([]bson.M, error) {
	coll := db.Client.Database("sample_db").Collection(collectionName)
	cursor, err := coll.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}
	var results []bson.M
	if err := cursor.All(context.TODO(), &results); err != nil {
		return nil, err
	}
	return results, nil
}

func (db *Database) Update(collectionName string, filter bson.D, update bson.D) (*mongo.UpdateResult, error) {
	coll := db.Client.Database("sample_db").Collection(collectionName)
	result, err := coll.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (db *Database) Delete(collectionName string, filter bson.D) (*mongo.DeleteResult, error) {
	coll := db.Client.Database("sample_db").Collection(collectionName)
	result, err := coll.DeleteMany(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	return result, nil
}
