#!/bin/bash
sam build

sam deploy --stack-name namik-chat-app --region us-west-2 --parameter-overrides 'TableName=simplechat_connections' --capabilities CAPABILITY_NAMED_IAM --on-failure ROLLBACK  --config-env default --confirm-changeset