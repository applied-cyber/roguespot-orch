package orchestrator

import (
	"context"
	"log"
	"net/http"
	"time"

	"roguespot-orch/alert"

	"github.com/gin-gonic/gin"
)

type AP struct {
	Address  string  `json:"address"`
	SSID     string  `json:"ssid"`
	Strength float64 `json:"strength"`
}

type Orchestrator struct {
	db       *DataStore
	alerts   *alert.AlertSys
	user     string
	password string
}

// Creates a new Orchestrator instance
func NewOrchestrator(ctx context.Context, uri, user, password string,
	alerts *alert.AlertSys) (*Orchestrator, error) {
	db, err := NewDataStore(ctx, uri, "roguespot-orch", "post-requests")
	if err != nil {
		return nil, err
	}
	return &Orchestrator{db: db,
		user: user, password: password,
		alerts: alerts}, nil
}

// Runs the orchestrator server
func (o *Orchestrator) Run() {
	authAccounts := gin.BasicAuth(gin.Accounts{
		o.user: o.password,
	})

	router := gin.Default()
	router.POST("/log", authAccounts, o.handlePost)
	o.alerts.Notify("Started orchestrator at " + time.Now().String())
	if err := router.Run(); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}

// Handles requests from the wardrivers
func (o *Orchestrator) handlePost(c *gin.Context) {
	var accessPoints []AP
	if err := c.BindJSON(&accessPoints); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for _, accessPoint := range accessPoints {
		err := o.handleAccessPoint(accessPoint, c)
		// TODO: Continue on error, and send back which access points had an error
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"response": "Access points handled correctly"})
}

func (o *Orchestrator) handleAccessPoint(accessPoint AP, c *gin.Context) error {
	log.Printf("Handling access point: %+v", accessPoint)
	ctx := c.Request.Context()

	// Check if the access point already exists
	exists, err := o.db.Exists(ctx, "address", accessPoint.Address)
	if err != nil {
		log.Printf("Error while checking if access point exists in database: %s", err)
		return err
	}

	if exists {
		log.Printf("Access point already exists in database. Skipping insert")
		return nil
	}

	// Insert the access point into the database
	err = o.db.Insert(ctx, accessPoint)
	if err != nil {
		log.Printf("Error while inserting access point into database: %s", err)
		return err
	}

	o.alerts.Notify("New access point detected: " + accessPoint.SSID + ", " + accessPoint.Address)

	return nil
}
