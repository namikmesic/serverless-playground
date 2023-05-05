package main

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/namikmesic/serverless-playground/src/lib/dynamo"
)

var (
	connectionHelper *dynamo.ConnectionHelper
)

func init() {
	tableName := os.Getenv("TABLE_NAME")
	connectionHelper = dynamo.NewConnectionHelper(tableName)
}

func handler(ctx context.Context, event events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Println("Connection established")

	// Store the connection ID in DynamoDB
	err := connectionHelper.WriteConnection(ctx, event.RequestContext.ConnectionID)
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
