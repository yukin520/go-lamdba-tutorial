package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func lamdaHandler(context context.Context, event json.RawMessage) (events.APIGatewayProxyResponse, error) {
	log.Println("Init lamdba handler.")

	response := events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       "\"Hello from Lambda!\"",
	}
	return response, nil
}

func main() {
	lambda.Start(lamdaHandler)
}
