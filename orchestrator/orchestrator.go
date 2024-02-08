package orchestrator

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Orchestrator struct {
	db DataStore
}

func NewOrchestrator(db DataStore) *Orchestrator {
	return &Orchestrator{db: db}
}

func (o *Orchestrator) Run() {
	router := gin.Default()
	router.POST("/log", o.handlePost)
	if err := router.Run(); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}

func (o *Orchestrator) handlePost(c *gin.Context) {
	var accessPoint AP
	if err := c.BindJSON(&accessPoint); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	containsAP, err := o.db.CheckIfAPExists(accessPoint)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if containsAP {
		c.JSON(http.StatusOK, gin.H{"response": "Access point already exists in database"})
		return
	}
	statusCode, responseBody := o.db.InsertAP(accessPoint)
	c.JSON(statusCode, gin.H{"response": responseBody})
}
