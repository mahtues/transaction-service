package transaction

import (
	"testing"

	"github.com/pkg/errors"
)

type mockRepository struct {
	saveTransactionFunc func(Transaction) error
	loadTransactionFunc func(id string) (Transaction, error)
}

func (m *mockRepository) SaveTransaction(transaction Transaction) error {
	return m.saveTransactionFunc(transaction)
}

func (m *mockRepository) LoadTransaction(id string) (Transaction, error) {
	return m.loadTransactionFunc(id)
}

func newServiceWithMockRepository(mock *mockRepository) Service {
	return Service{
		repository: mock,
	}
}

func TestUnit(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		s := newServiceWithMockRepository(&mockRepository{
			saveTransactionFunc: func(transaction Transaction) error {
				if transaction.Id == "" {
					return errors.New("transaction id not initialized")
				}

				return nil
			},
		})

		actual, actualErr := s.CreateTransaction(CreateRequest{})

		if actualErr != nil {
			t.Errorf("expected: %v, actual: %v", nil, actualErr)
		}

		if actual.Id == "" {
			t.Errorf("expected: %v, actual: %v", "random id", actual)
		}
	})

	t.Run("repository save error", func(t *testing.T) {
		s := newServiceWithMockRepository(&mockRepository{
			saveTransactionFunc: func(t Transaction) error {
				return errors.New("save error")
			},
		})

		_, actualErr := s.CreateTransaction(CreateRequest{})

		if actualErr == nil {
			t.Errorf("expected: %v, actual: %v", "not nil", actualErr)
		}
	})
}
