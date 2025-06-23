package handler

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"strconv"
	"strings"
	"time"

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

type todoRequestParams struct {
	Id          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}

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
		// Create Todo Item.
		if request.RequestContext.HTTP.Path == "/item" {
			var params todoRequestParams
			if len(request.Body) <= 0 {
				return infra.APIResponse(400, domain.ErrInvalidParameters)
			}
			err := json.Unmarshal([]byte(request.Body), &params)
			if err != nil {
				log.Printf("Unmarshal error: %v\n", err)
				return infra.APIResponse(400, domain.ErrInvalidParameters)
			}

			loc, _ := time.LoadLocation("Asia/Tokyo")
			todoData := domain.ToDo{
				Id:          params.Id,
				Name:        params.Name,
				Description: params.Description,
				Completed:   params.Completed,
				CreatedAt:   time.Now().In(loc),
				UpdatedAt:   time.Now().In(loc),
			}

			createdItemId, err := todoUsecase.CreateTodo(context, &todoData)
			if errors.Is(err, domain.ErrAlreadyExists) {
				body := struct{ Msg string }{
					Msg: "\"todo item is already exists.\"",
				}
				response, _ = infra.APIResponse(400, body)
				return response, nil
			}
			if err != nil {
				body := struct{ Msg string }{
					Msg: "\"faild to create todo item.\"",
				}
				response, _ = infra.APIResponse(500, body)
				return response, nil
			}

			body := struct {
				Msg string
				Id  uint
			}{
				Msg: "\"todo item is created.\"",
				Id:  createdItemId,
			}
			response, _ = infra.APIResponse(200, body)
			return response, nil
		}

		// Incorrect path.
		body := struct{ Msg string }{
			Msg: "\"not found.\"",
		}
		response, _ = infra.APIResponse(404, body)
		return response, nil
	}

	if string(request.RequestContext.HTTP.Method) == "PUT" {
		// Update Todo Item.
		if request.RequestContext.HTTP.Path == "/item" {
			var params todoRequestParams
			if len(request.Body) <= 0 {
				return infra.APIResponse(400, domain.ErrInvalidParameters)
			}
			err := json.Unmarshal([]byte(request.Body), &params)
			if err != nil {
				log.Printf("Unmarshal error: %v\n", err)
				return infra.APIResponse(400, domain.ErrInvalidParameters)
			}

			loc, _ := time.LoadLocation("Asia/Tokyo")
			todoData := domain.ToDo{
				Id:          params.Id,
				Name:        params.Name,
				Description: params.Description,
				Completed:   params.Completed,
				CreatedAt:   time.Now().In(loc),
				UpdatedAt:   time.Now().In(loc),
			}

			updatedTodoItem, err := todoUsecase.UpdateTodo(context, &todoData)
			if errors.Is(err, domain.ErrNotFound) {
				body := struct{ Msg string }{
					Msg: "\"todo item is not found.\"",
				}
				response, _ = infra.APIResponse(400, body)
				return response, nil
			}
			if err != nil {
				body := struct{ Msg string }{
					Msg: "\"faild to update todo item.\"",
				}
				response, _ = infra.APIResponse(500, body)
				return response, nil
			}
			response, _ = infra.APIResponse(200, updatedTodoItem)
			return response, nil
		}

		// Incorrect path.
		body := struct{ Msg string }{
			Msg: "\"not found.\"",
		}
		response, _ = infra.APIResponse(404, body)
		return response, nil
	}

	if string(request.RequestContext.HTTP.Method) == "DELETE" {
		// Delete Todo Item.
		if strings.HasPrefix(request.RawPath, "/item/") {
			todoItemId, err := strconv.Atoi(request.PathParameters["id"])
			if err != nil {
				body := struct{ Msg string }{
					Msg: "\"faild to delete todo item.\"",
				}
				response, _ = infra.APIResponse(500, body)
				return response, nil
			}

			err = todoUsecase.DeleteTodo(context, uint(todoItemId))
			if errors.Is(err, domain.ErrNotFound) {
				body := struct{ Msg string }{
					Msg: "\"todo item is not found.\"",
				}
				response, _ = infra.APIResponse(400, body)
				return response, nil
			}
			if err != nil {
				body := struct{ Msg string }{
					Msg: "\"faild to delete todo item.\"",
				}
				response, _ = infra.APIResponse(500, body)
				return response, nil
			}

			body := struct{ Msg string }{
				Msg: "\"todo item is delted sccessfully.\"",
			}
			response, _ = infra.APIResponse(200, body)
			return response, nil
		}

		// Incorrect path.
		body := struct{ Msg string }{
			Msg: "\"not found.\"",
		}
		response, _ = infra.APIResponse(404, body)
		return response, nil
	}

	body := struct{ Msg string }{
		Msg: "\"Http method is not supported.\"",
	}
	response, _ = infra.APIResponse(500, body)
	return response, errors.New("http method is not supported")
}
