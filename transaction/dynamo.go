package transaction

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type Dynamo struct {
	client *dynamodb.DynamoDB
}

func NewDynamoRepository(client *dynamodb.DynamoDB) *Dynamo {
	return &Dynamo{
		client: client,
	}
}

func (d *Dynamo) SaveTransaction() error {
	return nil
}

func (d *Dynamo) LoadTransaction() error {
	return nil
}
