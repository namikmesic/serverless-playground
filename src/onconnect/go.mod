module onconnect

go 1.20

replace github.com/namikmesic/serverless-playground/src/lib/dynamo => ../lib/dynamo

require (
	github.com/aws/aws-lambda-go v1.41.0
	github.com/namikmesic/serverless-playground/src/lib/dynamo v0.0.0-00010101000000-000000000000
)

require (
	github.com/aws/aws-sdk-go v1.44.257 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
)
