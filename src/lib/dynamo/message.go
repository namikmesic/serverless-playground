package dynamo

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type MessageHelper struct {
	DynamoHelper
}

func NewMessageHelper(tableName string) *MessageHelper {
	return &MessageHelper{
		DynamoHelper: *NewDynamoHelper(tableName),
	}
}

type Message struct {
	MessageID string `dynamodbav:"messageId"`
	Text      string `dynamodbav:"text"`
}

func (mh *MessageHelper) WriteMessage(ctx context.Context, messageID, text string) error {
	message := Message{
		MessageID: messageID,
		Text:      text,
	}
	return mh.WriteItem(ctx, message)
}

func (mh *MessageHelper) DeleteMessage(ctx context.Context, messageID string) error {
	key := map[string]*dynamodb.AttributeValue{
		"messageId": {
			S: aws.String(messageID),
		},
	}
	return mh.DeleteItem(ctx, key)
}

func (mh *MessageHelper) ScanMessages(ctx context.Context) (*dynamodb.ScanOutput, error) {
	return mh.ScanItems(ctx, "messageId, text")
}
