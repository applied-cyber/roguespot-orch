package orchestrator

import (
	"context"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDBHandler struct {
	collection *mongo.Collection
}

func NewMongoDBHandler(uri, dbName, collName string) *MongoDBHandler {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}
	collection := client.Database(dbName).Collection(collName)
	return &MongoDBHandler{collection: collection}
}

func (m *MongoDBHandler) CheckIfAPExists(accessPoint AP) (bool, error) {
	cursor, err := m.collection.Find(context.TODO(), bson.D{{Key: "address", Value: accessPoint.Address}})
	if err != nil {
		return false, err
	}
	defer cursor.Close(context.TODO())
	return cursor.Next(context.TODO()), nil
}

func (m *MongoDBHandler) InsertAP(accessPoint AP) (int, string) {
	insertResult, err := m.collection.InsertOne(context.Background(), accessPoint)
	if err != nil {
		return http.StatusInternalServerError, err.Error()
	}

	log.Printf("Insert successful, ID: %v", insertResult.InsertedID)
	return http.StatusCreated, "Successfully inserted access point"
}
