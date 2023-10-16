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

func (this Transaction) Equal(that Transaction) bool {
	e := this.Id == that.Id
	e = e && this.Description == that.Description
	e = e && this.Date.Equal(that.Date)
	e = e && this.AmountUs == that.AmountUs
	return e
}
