package main

import (
	"context"
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
	err := connectionHelper.DeleteConnection(ctx, event.RequestContext.ConnectionID)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: "Failed to disconnect: " + err.Error()}, err
	}

	return events.APIGatewayProxyResponse{StatusCode: 200, Body: "Disconnected."}, nil
}

func main() {
	lambda.Start(handler)
}
