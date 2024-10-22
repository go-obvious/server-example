package main

import (
	"context"

	"github.com/go-obvious/env"
	"github.com/go-obvious/server"

	"github.com/go-obvious/server-example/handlers"
	"github.com/go-obvious/server-example/internal/build"
	"github.com/go-obvious/server-example/store"
)

func main() {

	// validate the mode was passed
	_ = env.MustGet("SERVER_MODE")

	// Create the data storage layer
	ctx := context.Background()
	storage := store.NewDynamoDBStore(ctx, env.MustGet("TABLE"))

	// start the service
	server.New(
		build.Version(),
		handlers.NewProductService("/products", storage),
		handlers.NewPingService("/ping"),
	).Run(ctx)
}
