package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var collection *mongo.Collection

type AP struct {
	SSID     string  `json:"ssid"`
	Address  string  `json:"address"`
	Strength float64 `json:"strength"`
}

func collectionContainsAP(collection *mongo.Collection, accessPoint AP) (bool, error) {
	// Find all documents with the same MAC address as the access point
	log.Printf("Finding all documents with MAC address %s in collection", accessPoint.Address)
	cursor, err := collection.Find(context.TODO(), bson.D{{Key: "address", Value: accessPoint.Address}})

	defer func() {
		// Ensure the connection is always closed
		if err := cursor.Close(context.TODO()); err != nil {
			log.Printf("Error closing connection: %v", err)
		}
	}()

	if err != nil {
		log.Printf("Error finding access point in collection: %v", err)
		return false, err
	}

	return cursor.Next(context.TODO()), nil
}

func insertIntoCollection(collection *mongo.Collection, accessPoint AP) (int, string) {
	log.Printf("Inserting %+v", accessPoint)
	insertResult, err := collection.InsertOne(context.Background(), accessPoint)
	if err != nil {
		log.Printf("Error inserting into collection: %v", err)
		return http.StatusInternalServerError, err.Error()
	}

	log.Printf("Insert successful, ID: %v", insertResult.InsertedID)
	return http.StatusCreated, "Successfully inserted access point"
}

func handleAP(collection *mongo.Collection, accessPoint AP) (int, string) {
	containsAP, err := collectionContainsAP(collection, accessPoint)
	if err != nil {
		return http.StatusInternalServerError, err.Error()
	}

	if containsAP {
		log.Printf("Collection already contains access point with MAC address %s", accessPoint.Address)
		return http.StatusOK, "Access point already exists in database"
	}

	return insertIntoCollection(collection, accessPoint)
}

func postAP(c *gin.Context) {
	var accessPoint AP

	if err := c.BindJSON(&accessPoint); err != nil {
		log.Printf("Error binding JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	statusCode, responseBody := handleAP(collection, accessPoint)
	c.IndentedJSON(statusCode, gin.H{"response": responseBody})
}

func getMongoCollection(uri, dbName, collName string) *mongo.Collection {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	return client.Database(dbName).Collection(collName)
}

func main() {
	collection = getMongoCollection("mongodb://localhost:27017", "roguespot-orch", "post_requests")

	router := gin.Default()
	router.POST("/log", postAP)
	// TODO: Handle errors from Run?
	_ = router.Run()
}
