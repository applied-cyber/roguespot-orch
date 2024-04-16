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
	if config.Password == "" {
		log.Print(
			"!!! Warning: Using an empty password. It is highly recommended to use a password for ",
			"added security !!!",
		)
	}

	mongoURL := fmt.Sprintf("mongodb://%s:%d", config.MongoHost, config.MongoPort)
	orchestrator, err := orchestrator.NewOrchestrator(ctx, mongoURL, config.User, config.Password)
	if err != nil {
		log.Fatalf("Failed to initialize MongoDB handler: %v", err)
	}
	orchestrator.Run()
}
