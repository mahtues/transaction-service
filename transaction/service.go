package transaction

import (
	"time"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/mahtues/transaction-service/conversion"
)

type CreateRequest struct {
	Description string
	Date        time.Time
	AmountUs    string
}

type CreateResponse struct {
	Id string
}

type Service struct {
	repository Repository
	conversion conversion.Service
}

func NewServiceWithDynamo(client *dynamodb.DynamoDB, tableName string) Service {
	return Service{
		repository: NewDynamoRepository(client, tableName),
	}
}

func (s Service) GetTransaction(id string) error {
	return nil
}

func (s Service) CreateTransaction() error {
	return nil
}
