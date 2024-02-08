package main

import (
	"context"
	"log"
	"roguespot-orch/orchestrator"
)

func main() {
	ctx := context.Background()
	mongoDBHandler, err := orchestrator.NewMongoDBHandler(ctx, "mongodb://localhost:27017", "roguespot-orch", "post_requests")
	if err != nil {
		log.Fatalf("Failed to initialize MongoDB handler: %v", err)
	}
	orchestrator := orchestrator.NewOrchestrator(mongoDBHandler)
	orchestrator.Run()
}
