package service

import (
	"context"

	"web-server/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Database is a wrapper for the MongoDB client
type Database struct {
	Client *mongo.Client
}

// Disconnect closes the database connection
func (db *Database) Disconnect() error {
	return db.Client.Disconnect(context.TODO())
}

// Create inserts a new user into the database
func (db *Database) Create(collectionName string, user model.User) (*mongo.InsertOneResult, error) {
	coll := db.Client.Database("sample_db").Collection(collectionName)
	result, err := coll.InsertOne(context.TODO(), user)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Read fetches all users from the database
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

// Update modifies a user in the database
func (db *Database) Update(collectionName string, filter bson.D, update bson.D) (*mongo.UpdateResult, error) {
	coll := db.Client.Database("sample_db").Collection(collectionName)
	result, err := coll.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Delete removes a user from the database
func (db *Database) Delete(collectionName string, filter bson.D) (*mongo.DeleteResult, error) {
	coll := db.Client.Database("sample_db").Collection(collectionName)
	result, err := coll.DeleteMany(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	return result, nil
}
