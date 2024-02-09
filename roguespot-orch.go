package main

import (
	"context"
	"fmt"
	"log"

	"roguespot-orch/orchestrator"
)

func main() {
	ctx := context.Background()
	config := NewConfig()
	mongoURL := fmt.Sprintf("mongodb://%s:%d", config.MongoHost, config.MongoPort)

	mongoDBHandler, err := orchestrator.NewMongoDBHandler(ctx, mongoURL, "roguespot-orch", "post_requests")
	if err != nil {
		log.Fatalf("Failed to initialize MongoDB handler: %v", err)
	}
	orchestrator := orchestrator.NewOrchestrator(mongoDBHandler)
	orchestrator.Run()
}
