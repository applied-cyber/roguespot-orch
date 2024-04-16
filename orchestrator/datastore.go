package orchestrator

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Implements a data storage mechanism using MongoDB
type DataStore struct {
	collection *mongo.Collection
}

// Connects to a MongoDB and returns a DataStore
func NewDataStore(ctx context.Context, uri, dbName, collName string) (*DataStore, error) {
	log.Printf("Connecting to database '%s' at %s", dbName, uri)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}
	collection := client.Database(dbName).Collection(collName)
	return &DataStore{collection: collection}, nil
}

// Checks if a document with the specified key and value exists
func (ds *DataStore) Exists(ctx context.Context, key string, value interface{}) (bool, error) {
	filter := bson.D{{Key: key, Value: value}}
	cursor, err := ds.collection.Find(ctx, filter)
	if err != nil {
		return false, err
	}
	defer cursor.Close(ctx)
	return cursor.Next(ctx), nil
}

// Inserts a new document into the collection
func (ds *DataStore) Insert(ctx context.Context, document interface{}) error {
	_, err := ds.collection.InsertOne(ctx, document)
	if err != nil {
		return err
	}
	log.Printf("Insert successful")
	return nil
}
