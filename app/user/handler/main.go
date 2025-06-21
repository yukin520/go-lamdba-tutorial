package handler

import (
	"context"
	"errors"
	"log"
	"strconv"

	"github.com/aws/aws-lambda-go/events"

	"github.com/yukin520/go-lamdba-tutorial/app/domain"
	"github.com/yukin520/go-lamdba-tutorial/app/infra"
	_todoRepo "github.com/yukin520/go-lamdba-tutorial/app/user/repository"
	_todoUsecase "github.com/yukin520/go-lamdba-tutorial/app/user/usecase"
)

var (
	todoRepo    domain.TodoRepository
	todoUsecase domain.TodoUsecase
)

func init() {
	con, err := infra.NewDynamoDBConnectionFromEnv()
	if err != nil {
		log.Fatalf("Failed to connect to DynamoDB: %v", err)
	}
	todoRepo = _todoRepo.NewRepository(
		con.TableName,
		con.QueryIndexName,
		*con.Client,
	)
	todoUsecase = _todoUsecase.NewUsecase(todoRepo)
}

func LamdaHandler(context context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {
	log.Println("Init lamdba handler.")
	var response events.APIGatewayProxyResponse

	if string(request.RequestContext.HTTP.Method) == "GET" {
		// List Items.
		if request.RequestContext.HTTP.Path == "/" {
			todos, err := todoUsecase.ListTodo(context)
			if err != nil {
				body := struct{ Msg string }{
					Msg: "\"faild to fech todo list.\"",
				}
				response, _ = infra.APIResponse(500, body)
			} else {
				response, _ = infra.APIResponse(200, todos)
			}
			return response, nil
		}
		// Get Item by ID.
		if request.RequestContext.HTTP.Path == "/item" {
			todoItemId, err := strconv.Atoi(request.QueryStringParameters["id"])
			if err != nil {
				body := struct{ Msg string }{
					Msg: "\"faild to fech todo item.\"",
				}
				response, _ = infra.APIResponse(500, body)
				return response, nil
			}

			todo, err := todoUsecase.GetTodo(context, uint(todoItemId))
			if errors.Is(err, domain.ErrNotFound) {
				body := struct{ Msg string }{
					Msg: "\"todo item not found.\"",
				}
				response, _ = infra.APIResponse(404, body)
				return response, nil
			}
			if err != nil {
				body := struct{ Msg string }{
					Msg: "\"faild to fech todo item.\"",
				}
				response, _ = infra.APIResponse(500, body)
				return response, nil
			}
			response, _ = infra.APIResponse(200, todo)
			return response, nil
		}
	}

	if string(request.RequestContext.HTTP.Method) == "POST" {
		body := struct{ Msg string }{
			Msg: "\"Craeted item.\"",
		}
		response, _ = infra.APIResponse(200, body)
		return response, nil
	}
	if string(request.RequestContext.HTTP.Method) == "PUT" {
		body := struct{ Msg string }{
			Msg: "\"Updated item.\"",
		}
		response, _ = infra.APIResponse(200, body)
		return response, nil
	}
	if string(request.RequestContext.HTTP.Method) == "DELETE" {
		body := struct{ Msg string }{
			Msg: "\"Deleted item.\"",
		}
		response, _ = infra.APIResponse(200, body)
		return response, nil
	}

	body := struct{ Msg string }{
		Msg: "\"Http method is not supported.\"",
	}
	response, _ = infra.APIResponse(500, body)
	return response, errors.New("http method is not supported")
}
