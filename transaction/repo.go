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

func (this Transaction) Equal(that Transaction) bool {
	e := this.Id == that.Id
	e = e && this.Description == that.Description
	e = e && this.Date.Equal(that.Date)
	e = e && this.AmountUs == that.AmountUs
	return e
}
