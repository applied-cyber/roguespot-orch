package orchestrator

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// The orchestrator
type Orchestrator struct {
	db DataStore
}

func NewOrchestrator(db DataStore) *Orchestrator {
	return &Orchestrator{db: db}
}

func (o *Orchestrator) Run() {
	router := gin.Default()
	router.POST("/log", o.postAP)
	if err := router.Run(); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}

func (o *Orchestrator) handleAP(accessPoint AP) (int, string) {
	containsAP, err := o.db.CheckIfAPExists(accessPoint)
	if err != nil {
		return http.StatusInternalServerError, err.Error()
	}
	if containsAP {
		return http.StatusOK, "Access point already exists in database"
	}
	return o.db.InsertAP(accessPoint)
}

func (o *Orchestrator) postAP(c *gin.Context) {
	var accessPoint AP
	if err := c.BindJSON(&accessPoint); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	statusCode, responseBody := o.handleAP(accessPoint)
	c.JSON(statusCode, gin.H{"response": responseBody})
}
