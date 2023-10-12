package transaction

import (
	"time"
)

type Transaction struct {
	Id          string    `dynamodbav:"id"`
	Description string    `dynamodbav:"description"`
	Date        time.Time `dynamodbav:"date"`
	AmountUs    string    `dynamodbav:"amountUs"`
}

type Repository interface {
	SaveTransaction(transaction Transaction) error
	LoadTransaction(id string) (Transaction, error)
}
