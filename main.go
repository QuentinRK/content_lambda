package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/quentinrankin/content_lambda/handlers"
)

func main() {
	handler := handlers.Init()
	lambda.Start(handler)
}
