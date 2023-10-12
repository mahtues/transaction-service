package transaction

import (
	"testing"
)

func TestUnit(t *testing.T) {
	s := NewServiceWithDynamo(nil, "")

	t.Run("success", func(t *testing.T) {
		var expected, actual error = nil, s.CreateTransaction()

		if actual != expected {
			t.Errorf("expected: %v, actual: %v", expected, actual)
		}
	})

}
