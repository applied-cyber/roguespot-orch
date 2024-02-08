package orchestrator

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DataStore interface {
	CheckIfAPExists(ctx context.Context, ap AP) (bool, error)
	InsertAP(ctx context.Context, ap AP) error
}

type MongoDBHandler struct {
	collection *mongo.Collection
}

func NewMongoDBHandler(ctx context.Context, uri, dbName, collName string) (*MongoDBHandler, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}
	collection := client.Database(dbName).Collection(collName)
	return &MongoDBHandler{collection: collection}, nil
}

func (m *MongoDBHandler) CheckIfAPExists(ctx context.Context, accessPoint AP) (bool, error) {
	cursor, err := m.collection.Find(ctx, bson.D{{Key: "address", Value: accessPoint.Address}})
	if err != nil {
		return false, err
	}
	defer cursor.Close(ctx)
	return cursor.Next(ctx), nil
}

func (m *MongoDBHandler) InsertAP(ctx context.Context, accessPoint AP) error {
	_, err := m.collection.InsertOne(ctx, accessPoint)
	if err != nil {
		return err
	}
	log.Printf("Insert successful")
	return nil
}
