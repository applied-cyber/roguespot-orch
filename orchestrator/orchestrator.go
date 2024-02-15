package orchestrator

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Orchestrator struct {
	db       *DataStore
	user     string
	password string
}

func NewOrchestrator(ctx context.Context, uri, dbName, collName, user, password string) (*Orchestrator, error) {
	db, err := NewDataStore(ctx, uri, dbName, collName)
	orchestrator := &Orchestrator{db: db, user: user, password: password}
	return orchestrator, err
}

func (o *Orchestrator) Run() {
	authAccounts := gin.BasicAuth(gin.Accounts{
		o.user: o.password,
	})

	router := gin.Default()
	router.POST("/log", authAccounts, o.handlePost)
	if err := router.Run(); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}

// Handle a POST request from the wardriver
func (o *Orchestrator) handlePost(c *gin.Context) {
	var accessPoint AP
	if err := c.BindJSON(&accessPoint); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx := c.Request.Context()
	// Check if we already have an exact same record
	containsAP, err := o.db.APExists(ctx, accessPoint)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if containsAP {
		c.JSON(http.StatusOK, gin.H{"response": "Access point already exists in database"})
		return
	}
	// Insert the access point info into the database
	err = o.db.InsertAP(ctx, accessPoint)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"response": "Successfully inserted access point"})

	// TODO: publish to topic
}
