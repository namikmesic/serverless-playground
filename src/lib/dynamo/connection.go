package dynamo

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type ConnectionHelper struct {
	DynamoHelper
}

func NewConnectionHelper(tableName string) *ConnectionHelper {
	return &ConnectionHelper{
		DynamoHelper: *NewDynamoHelper(tableName),
	}
}

type Connection struct {
	ConnectionID string `dynamodbav:"connectionId"`
}

func (ch *ConnectionHelper) WriteConnection(ctx context.Context, connectionID string) error {
	connection := Connection{
		ConnectionID: connectionID,
	}
	return ch.WriteItem(ctx, connection)
}

func (ch *ConnectionHelper) DeleteConnection(ctx context.Context, connectionID string) error {
	key := map[string]*dynamodb.AttributeValue{
		"connectionId": {
			S: aws.String(connectionID),
		},
	}
	return ch.DeleteItem(ctx, key)
}

func (ch *ConnectionHelper) ScanConnections() (*dynamodb.ScanOutput, error) {
	return ch.ScanItems("connectionId")
}
