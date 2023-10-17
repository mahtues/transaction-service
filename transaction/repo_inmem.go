package transaction

import (
	"github.com/mahtues/transaction-service/apperrors"
	"github.com/pkg/errors"
)

type InMemRepository struct {
	mem map[string]Transaction
}

func (m *InMemRepository) Init() {
	m.mem = map[string]Transaction{}
}

func (m *InMemRepository) SaveTransaction(transaction Transaction) error {
	m.mem[transaction.Id] = transaction
	return nil
}

func (m *InMemRepository) LoadTransaction(id string) (Transaction, error) {
	transaction, found := m.mem[id]
	if !found {
		return Transaction{}, apperrors.Wrapf(errors.New("root cause"), "transaction id %s not found", id)
	}

	return transaction, nil
}
