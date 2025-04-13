package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/quentinrankin/content_lambda/handlers"
)

func main() {
	lambdaMaxRuntime := time.Now().Add(15 * time.Minute)

	ctx, cancel := context.WithDeadline(context.Background(), lambdaMaxRuntime)
	defer cancel()

	handler := handlers.Init()

	event := events.APIGatewayV2HTTPRequest{
		RequestContext: events.APIGatewayV2HTTPRequestContext{
			HTTP: events.APIGatewayV2HTTPRequestContextHTTPDescription{
				Path: "/projects",
			},
		},
	}

	res, err := handler(ctx, event)

	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Print(res.Body)

	os.Exit(0)
}
