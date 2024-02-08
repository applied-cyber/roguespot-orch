package main

import (
	"roguespot-orch/orchestrator"
)

func main() {
	mongoDBHandler := orchestrator.NewMongoDBHandler("mongodb://localhost:27017", "roguespot-orch", "post_requests")
	orchestrator := orchestrator.NewOrchestrator(mongoDBHandler)
	orchestrator.Run()
}
