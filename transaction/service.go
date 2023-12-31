package transaction

import (
	"time"

	"github.com/pkg/errors"

	"github.com/mahtues/transaction-service/conversion"
	"github.com/mahtues/transaction-service/random"
)

type Service struct {
	repository Repository
	conversion *conversion.Service
}

func NewService(repository Repository, conversion *conversion.Service) *Service {
	s := &Service{}
	s.Init(repository, conversion)
	return s
}

func (s *Service) Init(repository Repository, conversion *conversion.Service) {
	s.repository = repository
	s.conversion = conversion
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

	rate, err := s.conversion.GetRate(request.Country, transaction.Date)
	if err != nil {
		return GetResponse{}, errors.Wrap(err, "failed to get conversion rate")
	}

	amountConverted, err := conversion.Convert(transaction.AmountUs, rate)
	if err != nil {
		return GetResponse{}, errors.Wrap(err, "failed to convert amount")
	}

	return GetResponse{
		Id:              transaction.Id,
		Description:     transaction.Description,
		Date:            transaction.Date,
		AmountUs:        transaction.AmountUs,
		Rate:            rate,
		AmountConverted: amountConverted,
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
	Id      string
	Country string
}

type GetResponse struct {
	Id              string
	Description     string
	Date            time.Time
	AmountUs        string
	Rate            string
	AmountConverted string
}
