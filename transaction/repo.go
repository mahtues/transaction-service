package transaction

import (
	"time"
)

type Transaction struct {
	Id          string
	Description string
	Date        time.Time
	AmountUs    string
}

type Repository interface {
	SaveTransaction(transaction Transaction) error
	LoadTransaction(id string) (Transaction, error)
}
