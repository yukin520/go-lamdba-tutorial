package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/yukin520/go-lamdba-tutorial/app/user/handler"
)

func main() {
	lambda.Start(handler.LamdaHandler)
}
