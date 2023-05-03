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
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type Connection struct {
	ConnectionID string
}

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

func handler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	input := &dynamodb.ScanInput{
		ProjectionExpression: aws.String("connectionId"),
		TableName:            tableName,
	}

	connectionData, err := dynamodbSvc.Scan(input)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: err.Error()}, nil
	}

	endpoint := fmt.Sprintf("https://%s/%s", event.RequestContext.DomainName, event.RequestContext.Stage)
	apigwManagementAPI := apigatewaymanagementapi.New(session.Must(session.NewSession()), aws.NewConfig().WithEndpoint(endpoint))

	var postData map[string]interface{}
	err = json.Unmarshal([]byte(event.Body), &postData)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: err.Error()}, nil
	}

	data := postData["data"].(string)

	for _, item := range connectionData.Items {
		connection := Connection{}
		err := dynamodbattribute.UnmarshalMap(item, &connection)
		if err != nil {
			return events.APIGatewayProxyResponse{StatusCode: 500, Body: err.Error()}, nil
		}

		_, err = apigwManagementAPI.PostToConnection(&apigatewaymanagementapi.PostToConnectionInput{
			ConnectionId: aws.String(connection.ConnectionID),
			Data:         []byte(data),
		})

		if err != nil {
			_, ok := err.(*apigatewaymanagementapi.GoneException)
			if ok {
				fmt.Printf("Found stale connection, deleting %s\n", connection.ConnectionID)
				_, _ = dynamodbSvc.DeleteItem(&dynamodb.DeleteItemInput{
					Key:       map[string]*dynamodb.AttributeValue{"connectionId": {S: aws.String(connection.ConnectionID)}},
					TableName: tableName,
				})
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
