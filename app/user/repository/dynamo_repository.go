package repository

import (
	"context"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/yukin520/go-lamdba-tutorial/app/domain"
)

type repository struct {
	Table      *string
	QueryIndex *string
	Client     *dynamodb.Client
}

type TodoScan struct {
	Id          uint      `dynamodbav:"id"`
	Name        string    `dynamodbav:"name"`
	Description string    `dynamodbav:"description"`
	CreatedAt   time.Time `dynamodbav:"created_at"`
	UpdatedAt   time.Time `dynamodbav:"updated_at"`
	Completed   bool      `dynamodbav:"completed"`
}

func FromScan(s *TodoScan) *domain.ToDo {
	var todo domain.ToDo

	todo.Completed = s.Completed
	todo.CreatedAt = s.CreatedAt
	todo.Description = s.Description
	todo.Id = s.Id
	todo.Name = s.Name
	todo.UpdatedAt = s.UpdatedAt

	return &todo
}

func ToScan(t *domain.ToDo) *TodoScan {
	var todoScan TodoScan

	todoScan.Completed = t.Completed
	todoScan.CreatedAt = t.CreatedAt
	todoScan.Description = t.Description
	todoScan.Id = t.Id
	todoScan.Name = t.Name
	todoScan.UpdatedAt = t.UpdatedAt

	return &todoScan
}

func NewRepository(table string, query string, client dynamodb.Client) domain.TodoRepository {
	return &repository{
		Table:      &table,
		QueryIndex: &query,
		Client:     &client,
	}
}

func (m *repository) ListTodo(ctx context.Context) ([]*domain.ToDo, error) {
	RECORD_TYPE_KEY := "record_type"
	RECORD_TYPE_VALUE := "todo"

	var todoList []*domain.ToDo

	keyEx := expression.Key(RECORD_TYPE_KEY).Equal(expression.Value(RECORD_TYPE_VALUE))
	expr, err := expression.NewBuilder().WithKeyCondition(keyEx).Build()

	if err != nil {
		log.Printf("Couldn't build expression for query. Here's why: %v\n", err)
		return []*domain.ToDo{}, err
	}

	queryPaginator := dynamodb.NewQueryPaginator(m.Client, &dynamodb.QueryInput{
		TableName:                 aws.String(*m.Table),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
		IndexName:                 aws.String(*m.QueryIndex),
	})
	for queryPaginator.HasMorePages() {
		var todoPages []TodoScan

		response, err := queryPaginator.NextPage(ctx)
		if err != nil {
			log.Printf("Couldn't query for todo. Here's why: %v\n", err)
			return nil, err
		}

		err = attributevalue.UnmarshalListOfMaps(response.Items, &todoPages)
		if err != nil {
			log.Printf("Couldn't unmarshal query response. Here's why: %v\n", err)
			return nil, err
		}

		for _, v := range todoPages {
			todoList = append(todoList, FromScan(&v))
		}
	}
	return todoList, nil
}

func (m *repository) GetTodo(ctx context.Context, id uint) (*domain.ToDo, error) {
	panic("GetTodo not implemented")
}

func (m *repository) CreateTodo(ctx context.Context, param *domain.ToDo) (uint, error) {
	panic("GetTodo not implemented")
}

func (m *repository) UpdateTodo(ctx context.Context, param *domain.ToDo) (*domain.ToDo, error) {
	panic("GetTodo not implemented")
}

func (m *repository) DeleteTodo(ctx context.Context, id uint) error {
	panic("GetTodo not implemented")
}
