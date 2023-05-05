module sendmessage

go 1.20

require (
	github.com/aws/aws-lambda-go v1.40.0
	github.com/aws/aws-sdk-go v1.44.257
	github.com/namikmesic/serverless-playground/src/lib/dynamo v0.0.0-00010101000000-000000000000
)

require github.com/jmespath/go-jmespath v0.4.0 // indirect

replace github.com/namikmesic/serverless-playground/src/lib/dynamo => ../lib/dynamo
