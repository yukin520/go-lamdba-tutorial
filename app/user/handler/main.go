package handler

import (
	"context"
	"errors"
	"log"

	"github.com/aws/aws-lambda-go/events"

	"github.com/yukin520/go-lamdba-tutorial/app/infra"
)

func LamdaHandler(context context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {
	log.Println("Init lamdba handler.")
	var response events.APIGatewayProxyResponse

	if string(request.RequestContext.HTTP.Method) == "GET" {
		body := struct{ Name string }{
			Name: "hoge",
		}
		response, _ = infra.APIResponse(200, body)
	} else if string(request.RequestContext.HTTP.Method) == "POST" {
		body := struct{ Msg string }{
			Msg: "\"Craeted item.\"",
		}
		response, _ = infra.APIResponse(200, body)
	} else if string(request.RequestContext.HTTP.Method) == "PUT" {
		body := struct{ Msg string }{
			Msg: "\"Updated item.\"",
		}
		response, _ = infra.APIResponse(200, body)
	} else if string(request.RequestContext.HTTP.Method) == "DELETE" {
		body := struct{ Msg string }{
			Msg: "\"Deleted item.\"",
		}
		response, _ = infra.APIResponse(200, body)
	} else {
		body := struct{ Msg string }{
			Msg: "\"Http method is not supported.\"",
		}
		response, _ = infra.APIResponse(200, body)
		return response, errors.New("http method is not supported")
	}

	return response, nil

}
