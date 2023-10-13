package support

import (
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func AwsSessionLocalhost(port int) (*session.Session, error) {
	return session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Region:      aws.String("us-east-2"),
			Endpoint:    aws.String(fmt.Sprintf("http://localhost:%d", port)),
			Credentials: credentials.NewStaticCredentials("fake-access-key", "fake-secret-key", "fake-token"),
		},
		SharedConfigState: session.SharedConfigEnable,
	})
}

func AwsDynamoDbClient(session *session.Session) (*dynamodb.DynamoDB, error) {
	client := dynamodb.New(session)
	if client == nil {
		return nil, errors.New("error creating dynamodb client")
	}

	return client, nil
}
