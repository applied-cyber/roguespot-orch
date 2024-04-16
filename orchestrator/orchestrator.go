package orchestrator

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AP struct {
	Address  string `json:"address"`
	SSID     string `json:"ssid"`
	Security string `json:"security"`
}

type Orchestrator struct {
	db       *DataStore
	user     string
	password string
}

// Creates a new Orchestrator instance
func NewOrchestrator(ctx context.Context, uri, user, password string) (*Orchestrator, error) {
	db, err := NewDataStore(ctx, uri, "roguespot-orch", "post-requests")
	if err != nil {
		return nil, err
	}
	return &Orchestrator{db: db, user: user, password: password}, nil
}

// Runs the orchestrator server
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

// Handles requests from the wardrivers
func (o *Orchestrator) handlePost(c *gin.Context) {
	var accessPoint AP
	if err := c.BindJSON(&accessPoint); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx := c.Request.Context()
	// Check if the access point already exists
	exists, err := o.db.Exists(ctx, "address", accessPoint.Address)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if exists {
		c.JSON(http.StatusOK, gin.H{"response": "Access point already exists in database"})
		return
	}
	// Insert the access point into the database
	err = o.db.Insert(ctx, accessPoint)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"response": "Successfully inserted access point"})
}
