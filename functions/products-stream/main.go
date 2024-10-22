package main

import (
	"context"
	"os"

	"github.com/aws/aws-lambda-go/lambda"

	"github.com/go-obvious/server-example/bus"
	"github.com/go-obvious/server-example/domain"
	"github.com/go-obvious/server-example/handlers"
)

func main() {
	eventBusName, ok := os.LookupEnv("EVENT_BUS_NAME")
	if !ok {
		panic("Need EVENT_BUS_NAME environment variable")
	}

	store := bus.NewEventBridgeBus(context.TODO(), eventBusName)
	domain := domain.NewProductsStream(store)
	handler := handlers.NewDynamoDBEventHandler(domain)
	lambda.Start(handler.StreamHandler)
}
