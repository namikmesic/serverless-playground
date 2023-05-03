package main

import (
	"context"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

var (
	tableName   = aws.String(os.Getenv("TABLE_NAME"))
	dynamodbSvc *dynamodb.DynamoDB
)

func init() {
	// Create a new session and DynamoDB client
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	dynamodbSvc = dynamodb.New(sess)
}

func handler(ctx context.Context, event events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	deleteParams := &dynamodb.DeleteItemInput{
		TableName: tableName,
		Key: map[string]*dynamodb.AttributeValue{
			"connectionId": {
				S: aws.String(event.RequestContext.ConnectionID),
			},
		},
	}

	_, err := dynamodbSvc.DeleteItemWithContext(ctx, deleteParams)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: "Failed to disconnect: " + err.Error()}, err
	}

	return events.APIGatewayProxyResponse{StatusCode: 200, Body: "Disconnected."}, nil
}

func main() {
	lambda.Start(handler)
}
