package dynamo

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type DynamoOperations interface {
	WriteItem(ctx context.Context, item interface{}) error
	ReadItem(ctx context.Context, key map[string]*dynamodb.AttributeValue) (map[string]*dynamodb.AttributeValue, error)
	DeleteItem(ctx context.Context, key map[string]*dynamodb.AttributeValue) error
	ScanItems(projectionExpression string) (*dynamodb.ScanOutput, error)
}

type DynamoHelper struct {
	TableName   *string
	DynamoDBSvc *dynamodb.DynamoDB
}

func NewDynamoHelper(tableName string) *DynamoHelper {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	dynamoDBSvc := dynamodb.New(sess)

	return &DynamoHelper{
		TableName:   aws.String(tableName),
		DynamoDBSvc: dynamoDBSvc,
	}
}

func (dh *DynamoHelper) WriteItem(ctx context.Context, item interface{}) error {
	av, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		return err
	}

	putParams := &dynamodb.PutItemInput{
		TableName: dh.TableName,
		Item:      av,
	}

	_, err = dh.DynamoDBSvc.PutItemWithContext(ctx, putParams)
	return err
}

func (dh *DynamoHelper) ReadItem(ctx context.Context, key map[string]*dynamodb.AttributeValue) (map[string]*dynamodb.AttributeValue, error) {
	getParams := &dynamodb.GetItemInput{
		TableName: dh.TableName,
		Key:       key,
	}

	result, err := dh.DynamoDBSvc.GetItemWithContext(ctx, getParams)
	if err != nil {
		return nil, err
	}

	return result.Item, nil
}

func (dh *DynamoHelper) DeleteItem(ctx context.Context, key map[string]*dynamodb.AttributeValue) error {
	deleteParams := &dynamodb.DeleteItemInput{
		TableName: dh.TableName,
		Key:       key,
	}

	_, err := dh.DynamoDBSvc.DeleteItemWithContext(ctx, deleteParams)
	return err
}

func (dh *DynamoHelper) ScanItems(projectionExpression string) (*dynamodb.ScanOutput, error) {
	input := &dynamodb.ScanInput{
		ProjectionExpression: aws.String(projectionExpression),
		TableName:            dh.TableName,
	}

	return dh.DynamoDBSvc.Scan(input)
}
