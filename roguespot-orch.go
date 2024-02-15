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
	orchestrator, err := orchestrator.NewOrchestrator(ctx, mongoURL, "roguespot-orch", "post-requests")
	if err != nil {
		log.Fatalf("Failed to initialize MongoDB handler: %v", err)
	}
	orchestrator.Run()
}
