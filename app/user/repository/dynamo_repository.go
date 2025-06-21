package repository

import (
	"context"
	"errors"
	"log"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
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
	RecordType  string    `dynamodbav:"record_type"`
}

var (
	RECORD_TYPE_KEY   string = "record_type"
	RECORD_TYPE_VALUE string = "todo"
)

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
	todoScan.RecordType = "todo"

	return &todoScan
}

func NewRepository(table string, query string, client dynamodb.Client) domain.TodoRepository {
	return &repository{
		Table:      &table,
		QueryIndex: &query,
		Client:     &client,
	}
}

func (m *repository) getTodoItemById(ctx context.Context, id uint) (*TodoScan, error) {
	var todo TodoScan
	strId := strconv.FormatUint(uint64(id), 10)

	idMap := map[string]types.AttributeValue{"id": &types.AttributeValueMemberN{Value: strId}}

	response, err := m.Client.GetItem(ctx, &dynamodb.GetItemInput{
		Key: idMap, TableName: aws.String(*m.Table),
	})

	if err != nil {
		return nil, err
	}
	if response.Item == nil || len(response.Item) == 0 {
		return nil, domain.ErrNotFound
	}

	err = attributevalue.UnmarshalMap(response.Item, &todo)
	if err != nil {
		return nil, err
	}

	return &todo, nil
}

func (m *repository) ListTodo(ctx context.Context) ([]*domain.ToDo, error) {

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
	todoScan, err := m.getTodoItemById(ctx, id)
	if err != nil {
		log.Printf("Couldn't get todo item from dynamodb. Here's why: %v\n", err)
		return nil, err
	}
	return FromScan(todoScan), nil
}

func (m *repository) CreateTodo(ctx context.Context, param *domain.ToDo) (uint, error) {
	currentItem, err := m.getTodoItemById(ctx, param.Id)
	if currentItem != nil {
		return param.Id, domain.ErrAlreadyExists
	}
	// The absence of the item is a prerequisite for creating a new one.
	if err != nil && !errors.Is(err, domain.ErrNotFound) {
		log.Printf("Couldn't get todo item from dynamodb. Here's why: %v\n", err)
		return param.Id, err
	}

	// Create item record to DynamoDB talbe.
	todoScan := ToScan(param)
	item, err := attributevalue.MarshalMap(todoScan)
	if err != nil {
		return param.Id, err
	}
	_, err = m.Client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(*m.Table), Item: item,
	})
	if err != nil {
		log.Printf("Couldn't add item to table. Here's why: %v\n", err)
	}
	return param.Id, err
}

func (m *repository) UpdateTodo(ctx context.Context, param *domain.ToDo) (*domain.ToDo, error) {
	// fetch target todo item from DynamoDB.
	currentItem, err := m.getTodoItemById(ctx, param.Id)
	if err != nil {
		log.Printf("Couldn't get todo item from dynamodb. Here's why: %v\n", err)
		return nil, err
	}

	// Update todo item record at DynamoDB.
	var response *dynamodb.UpdateItemOutput
	update := expression.Set(expression.Name("name"), expression.Value(param.Name))
	update.Set(expression.Name("description"), expression.Value(param.Description))
	update.Set(expression.Name("completed"), expression.Value(param.Completed))
	update.Set(expression.Name("updated_at"), expression.Value(param.UpdatedAt))
	update.Set(expression.Name("created_at"), expression.Value(currentItem.CreatedAt))   // Use current item's value.
	update.Set(expression.Name("record_type"), expression.Value(currentItem.RecordType)) // Use current item's value.
	expr, err := expression.NewBuilder().WithUpdate(update).Build()
	if err != nil {
		log.Printf("Couldn't build expression for update. Here's why: %v\n", err)
		return nil, err
	} else {
		idMap := map[string]types.AttributeValue{"id": &types.AttributeValueMemberN{Value: strconv.FormatUint(uint64(currentItem.Id), 10)}}
		response, err = m.Client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
			TableName:                 aws.String(*m.Table),
			Key:                       idMap,
			ExpressionAttributeNames:  expr.Names(),
			ExpressionAttributeValues: expr.Values(),
			UpdateExpression:          expr.Update(),
			ReturnValues:              types.ReturnValueUpdatedNew,
		})
		if err != nil {
			log.Printf("Couldn't update todo %v. Here's why: %v\n", param.Id, err)
			return nil, err
		}
	}

	// perse dynamodb data to todo item.
	var todo TodoScan
	err = attributevalue.UnmarshalMap(response.Attributes, &todo)
	if err != nil {
		log.Printf("Couldn't unmarshall update response. Here's why: %v\n", err)
		return nil, err
	}
	return FromScan(&todo), nil
}

func (m *repository) DeleteTodo(ctx context.Context, id uint) error {
	panic("GetTodo not implemented")
}
