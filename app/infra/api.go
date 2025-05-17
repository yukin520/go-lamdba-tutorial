package infra

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
)

// APIResponse ...
func APIResponse(statusCode int, body interface{}) (events.APIGatewayProxyResponse, error) {
	bytes, _ := json.Marshal(&body)

	return events.APIGatewayProxyResponse{
		Body:       string(bytes),
		StatusCode: statusCode,
	}, nil
}
