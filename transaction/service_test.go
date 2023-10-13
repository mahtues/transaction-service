package transaction

import (
	"testing"

	"github.com/mahtues/transaction-service/conversion"
	"github.com/pkg/errors"
)

func TestUnit(t *testing.T) {
	base := NewService(nil, conversion.Service{})

	t.Run("success", func(t *testing.T) {
		s := newServiceWithMockRepository(base, &mockRepository{
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
		s := newServiceWithMockRepository(base, &mockRepository{
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

type mockRepository struct {
	saveTransactionFunc func(Transaction) error
	loadTransactionFunc func(id string) (Transaction, error)
	next                Repository
}

func newServiceWithMockRepository(service *Service, mock *mockRepository) *Service {
	newService := *service
	mock.next = service.repository
	newService.repository = mock

	return &newService
}

func (m *mockRepository) SaveTransaction(transaction Transaction) error {
	if m.saveTransactionFunc == nil {
		return m.next.SaveTransaction(transaction)
	}
	return m.saveTransactionFunc(transaction)
}

func (m *mockRepository) LoadTransaction(id string) (Transaction, error) {
	if m.loadTransactionFunc == nil {
		return m.next.LoadTransaction(id)
	}
	return m.loadTransactionFunc(id)
}
