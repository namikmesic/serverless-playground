package main

import (
	"context"
	"log"
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
	log.Println("Connection established")

	// Store the connection ID in DynamoDB
	putParams := &dynamodb.PutItemInput{
		TableName: tableName,
		Item: map[string]*dynamodb.AttributeValue{
			"connectionId": {
				S: aws.String(event.RequestContext.ConnectionID),
			},
		},
	}
	_, err := dynamodbSvc.PutItemWithContext(ctx, putParams)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: "Failed to connect: " + err.Error()}, nil
	}

	response := events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       "Connected.",
	}

	return response, nil
}

func main() {
	lambda.Start(handler)
}
