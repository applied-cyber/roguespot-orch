package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var collection *mongo.Collection

type AP struct {
	SSID string `json:"ssid"`
	MAC  string `json:"mac"`
}

func postAP(c *gin.Context) {
	var newAP AP

	if err := c.BindJSON(&newAP); err != nil {
		log.Printf("Error binding JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Attempting to insert: %+v", newAP)

	// Insert the newAP into the collection
	insertResult, err := collection.InsertOne(context.TODO(), newAP)
	if err != nil {
		log.Printf("Error inserting into collection: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Insert successful, ID: %v", insertResult.InsertedID)

	// c.IndentedJSON(http.StatusCreated, newAP)

	// Respond with the inserted document's ID
	c.IndentedJSON(http.StatusCreated, gin.H{"_id": insertResult.InsertedID})
}

func getMongoCollection(uri, dbName, collName string) *mongo.Collection {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	return client.Database(dbName).Collection(collName)
}

func main() {
	collection = getMongoCollection("mongodb://localhost:27017", "roguespot-orch", "post_requests")

	router := gin.Default()
	router.POST("/log", postAP)
	router.Run()
}
