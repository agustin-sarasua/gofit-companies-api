package main

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

const tableName = "Companies"

var companiesGSI = "companiesGSI"

// Declare a new DynamoDB instance. Note that this is safe for concurrent
// use.
var db dynamodbiface.DynamoDBAPI = dynamodb.New(session.New(), aws.NewConfig().WithRegion("us-east-1"))

func addType(av map[string]*dynamodb.AttributeValue, itype string) {
	av["DocType"] = &dynamodb.AttributeValue{
		S: aws.String(itype),
	}
}

func putCompany(e *Company) error {
	av, err := dynamodbattribute.MarshalMap(e)
	if err != nil {
		panic(fmt.Sprintf("failed to DynamoDB marshal Record, %v", err))
	}
	addType(av, "Company")
	e.Staff = nil

	input := &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      av,
	}

	_, err = db.PutItem(input)
	return err
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

func getUserCompanies(userSub string, limit int64) ([]*Company, error) {
	log.Printf("Loading user companies for %s", userSub)

	// Construct the Key condition builder
	keyCond := expression.Key("UserSub").Equal(expression.Value(userSub))

	// Construct the filter builder with a name and value.
	filt := expression.Name("DocType").Equal(expression.Value("Company"))

	// Using the filter and projections create a DynamoDB expression from the two.
	expr, err := expression.NewBuilder().
		WithKeyCondition(keyCond).
		WithFilter(filt).
		Build()
	if err != nil {
		fmt.Println(err)
	}

	// Prepare the input for the query.
	input := &dynamodb.QueryInput{
		TableName: aws.String(tableName),
		Limit:     &limit,
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
		FilterExpression:          expr.Filter(),
	}

	resp, err := db.Query(input)
	if err == nil {
		fmt.Println(resp)
		ps := []*Company{}
		err = dynamodbattribute.UnmarshalListOfMaps(resp.Items, &ps)
		log.Printf("Response count %d", *resp.Count)
		return ps, nil
	}
	log.Printf("Error %s", err.Error())
	return nil, err
}

func getCompanyWithStaff(companyID string, limit int64) (*Company, error) {

	// Prepare the input for the query.
	input := &dynamodb.QueryInput{
		TableName: aws.String(tableName),
		IndexName: &companiesGSI,
		Limit:     &limit,
		KeyConditions: map[string]*dynamodb.Condition{
			"CompanyID": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String(companyID),
					},
				},
			},
		},
	}
	var resp, err = db.Query(input)
	if err == nil {
		var company Company
		ss := []*Staff{}
		// err = dynamodbattribute.UnmarshalListOfMaps(resp.Items, &ps)
		for _, m := range resp.Items {
			if t, ok := m["DocType"]; ok {
				if t.GoString() == "Company" {
					err = dynamodbattribute.UnmarshalMap(m, &company)
				} else if t.GoString() == "Staff" {
					c := Staff{}
					err = dynamodbattribute.UnmarshalMap(m, &c)
					if err == nil {
						ss = append(ss, &c)
					}
				}
			}
		}
		company.Staff = ss
		return &company, nil
	}
	return nil, err

}
