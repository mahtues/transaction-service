package transaction

import (
	"time"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/pkg/errors"

	"github.com/mahtues/transaction-service/conversion"
	"github.com/mahtues/transaction-service/random"
)

type Service struct {
	repository Repository
	conversion conversion.Service
}

func NewServiceWithDynamo(client *dynamodb.DynamoDB, tableName string) Service {
	return Service{
		repository: NewDynamoRepository(client, tableName),
	}
}

func (s *Service) CreateTransaction(request CreateRequest) (CreateResponse, error) {
	transaction := Transaction{
		Id:          "tr-" + random.String(20),
		Description: request.Description,
		Date:        request.Date,
		AmountUs:    request.AmountUs,
	}

	if err := s.repository.SaveTransaction(transaction); err != nil {
		return CreateResponse{}, errors.Wrap(err, "failed to save transaction")
	}

	return CreateResponse{
		Id: transaction.Id,
	}, nil
}

func (s *Service) GetTransaction(request GetRequest) (GetResponse, error) {
	transaction, err := s.repository.LoadTransaction(request.Id)
	if err != nil {
		return GetResponse{}, errors.Wrapf(err, "failed to load transaction id %v", request.Id)
	}

	return GetResponse{
		Id:              transaction.Id,
		Description:     transaction.Description,
		Date:            transaction.Date,
		AmountUs:        transaction.AmountUs,
		Rate:            "",
		AmountConverted: "",
	}, nil
}

type CreateRequest struct {
	Description string
	Date        time.Time
	AmountUs    string
}

type CreateResponse struct {
	Id string
}

type GetRequest struct {
	Id string
}

type GetResponse struct {
	Id              string
	Description     string
	Date            time.Time
	AmountUs        string
	Rate            string
	AmountConverted string
}
