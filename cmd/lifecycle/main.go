package main

import (
	"context"

	"github.com/go-obvious/server"

	"github.com/go-obvious/server-example/internal/build"
	"github.com/go-obvious/server-example/internal/service/database"
	"github.com/go-obvious/server-example/internal/service/worker"
)

func main() {
	// Example demonstrating API lifecycle management with graceful shutdown
	// Shows how to use configuration registry and lifecycle hooks

	server.New(
		&server.ServerVersion{
			Revision: build.Rev,
			Tag:      build.Tag,
			Time:     build.Time,
		},
	).WithAPIs(
		database.NewService(), // Service with lifecycle hooks
		worker.NewService(),   // Background worker with lifecycle
	).Run(context.Background())

	// Lifecycle flow:
	// 1. Configuration loaded via registry
	// 2. database.Start() - connects to database
	// 3. worker.Start() - starts background tasks
	// 4. Server accepts requests
	// 5. On SIGTERM/SIGINT: graceful shutdown begins
	// 6. worker.Stop() - stops background tasks
	// 7. database.Stop() - closes database connections
	// 8. Server waits for existing requests to complete (30s timeout)
}
