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

func (ds *DataStore) APExists(ctx context.Context, accessPoint AP) (bool, error) {
	cursor, err := ds.collection.Find(ctx, bson.D{{Key: "address", Value: accessPoint.Address}})
	if err != nil {
		return false, err
	}
	defer cursor.Close(ctx)
	return cursor.Next(ctx), nil
}

func (ds *DataStore) InsertAP(ctx context.Context, accessPoint AP) error {
	_, err := ds.collection.InsertOne(ctx, accessPoint)
	if err != nil {
		return err
	}
	log.Printf("Insert successful")
	return nil
}
