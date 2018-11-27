package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

const tableName = "Companies"

var companiesUserSubGSI = "companiesUserSubGSI"

// Declare a new DynamoDB instance. Note that this is safe for concurrent
// use.
var db dynamodbiface.DynamoDBAPI = dynamodb.New(session.New(), aws.NewConfig().WithRegion("us-east-1"))

func addType(av map[string]*dynamodb.AttributeValue, itype string) {
	av["DocType"] = &dynamodb.AttributeValue{
		S: aws.String(itype),
	}
}

func putStaff(s *Staff) error {
	av, err := dynamodbattribute.MarshalMap(s)
	addType(av, "Staff")
	if err != nil {
		panic(fmt.Sprintf("failed to DynamoDB marshal Record, %v", err))
	}
	input := &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      av,
	}
	_, err = db.PutItem(input)
	return err
}