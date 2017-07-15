package dynamo_db

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type DynamoDb struct {
	db         *dynamodb.DynamoDB
	table_name *string
}

func NewDynamoDb(table_name string) *DynamoDb {
	db := dynamodb.New(session.New(&aws.Config{
		Region:      region,
		Credentials: cred,
	}))
	return &DynamoDb{
		db:         db,
		table_name: aws.String(table_name),
	}
}
