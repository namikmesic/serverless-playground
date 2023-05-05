package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/apigatewaymanagementapi"
	"github.com/namikmesic/serverless-playground/src/lib/dynamo" // Update the import path to your package location
)

var (
	connectionHelper *dynamo.ConnectionHelper
)

func init() {
	tableName := os.Getenv("TABLE_NAME")
	connectionHelper = dynamo.NewConnectionHelper(tableName)
}

func handler(ctx context.Context, event events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Scan for all connections
	connections, err := connectionHelper.ScanConnections(ctx)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: err.Error()}, nil
	}

	// Get the current connection ID
	currentConnectionID := event.RequestContext.ConnectionID

	endpoint := fmt.Sprintf("https://%s/%s", event.RequestContext.DomainName, event.RequestContext.Stage)
	apigwManagementAPI := apigatewaymanagementapi.New(session.Must(session.NewSession()), aws.NewConfig().WithEndpoint(endpoint))

	var postData map[string]interface{}
	err = json.Unmarshal([]byte(event.Body), &postData)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: err.Error()}, nil
	}

	data := postData["data"].(string)

	for _, connection := range connections {
		if connection == currentConnectionID {
			continue // Skip the current connection
		}

		_, err = apigwManagementAPI.PostToConnection(&apigatewaymanagementapi.PostToConnectionInput{
			ConnectionId: aws.String(connection),
			Data:         []byte(data),
		})

		if err != nil {
			_, ok := err.(*apigatewaymanagementapi.GoneException)
			if ok {
				fmt.Printf("Found stale connection, deleting %s\n", connection)
				_ = connectionHelper.DeleteConnection(ctx, connection)
			} else {
				return events.APIGatewayProxyResponse{StatusCode: 500, Body: err.Error()}, nil
			}
		}
	}

	return events.APIGatewayProxyResponse{StatusCode: 200, Body: "Data sent."}, nil
}
func main() {
	lambda.Start(handler)
}
