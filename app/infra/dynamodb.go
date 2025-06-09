package infra

import (
	"context"
	"errors"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type DynamoDBConnection struct {
	Client         *dynamodb.Client
	TableName      string
	QueryIndexName string
}

func NewDynamoDBConnectionFromEnv() (*DynamoDBConnection, error) {
	region, isExist := os.LookupEnv("AWS_REGION")
	if !isExist {
		// default region is "ap-northeast-1".
		region = "ap-northeast-1"
	}

	tableName, isExist := os.LookupEnv("DYNAMODB_TABLE")
	if !isExist {
		return nil, errors.New("can not be retrieved DyanmoDb table name from environment variable DYNAMODB_TABLE")
	}

	queryIndexName, isExist := os.LookupEnv("QUERY_INDEX")
	if !isExist {
		return nil, errors.New("can not be retrieved DyanmoDb table name from environment variable QUERY_INDEX")
	}

	accessKey, _ := os.LookupEnv("AWS_ACCESS_KEY_ID")
	secretAccessKey, _ := os.LookupEnv("AWS_SECRET_ACCESS_KEY")
	dynamodbEndpoint, _ := os.LookupEnv("DYNAMODB_ENDPOINT")

	return NewDynamoDBConnection(region, tableName, queryIndexName, accessKey, secretAccessKey, dynamodbEndpoint)
}

func NewDynamoDBConnection(region string, tableName string, queryIndexName string, accessKey string, secretAccessKey string, dynamodbEndpoint string) (*DynamoDBConnection, error) {
	var (
		cfg    aws.Config
		client *dynamodb.Client
		err    error
	)

	if isValidAccessKeyAndSecretAccessKey(accessKey, secretAccessKey) {
		credential := aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(accessKey, secretAccessKey, ""))
		cfg, err = config.LoadDefaultConfig(context.TODO(),
			config.WithRegion(region),
			config.WithCredentialsProvider(credential))
		if err != nil {
			return nil, err
		}
	} else {
		cfg, err = config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
		if err != nil {
			return nil, err
		}
	}

	if dynamodbEndpoint != "" {
		// execute with localstack.
		client = dynamodb.NewFromConfig(cfg, func(o *dynamodb.Options) {
			o.BaseEndpoint = aws.String(dynamodbEndpoint)
		})
	} else {
		client = dynamodb.NewFromConfig(cfg)
	}

	con := DynamoDBConnection{
		Client:         client,
		TableName:      tableName,
		QueryIndexName: queryIndexName,
	}
	return &con, nil
}

func isValidAccessKeyAndSecretAccessKey(accessKey string, secretAccessKey string) bool {
	if accessKey != "" && secretAccessKey != "" {
		return true
	}
	return false
}
