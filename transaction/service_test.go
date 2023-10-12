package transaction

import (
	"flag"
	"testing"
)

var integration = flag.Bool("integration", false, "only perform local tests")

func TestUnit(t *testing.T) {
	s := NewServiceWithDynamo(nil)

	t.Run("success", func(t *testing.T) {
		var expected, actual error = nil, s.CreateTransaction()

		if actual != expected {
			t.Errorf("expected: %v, actual: %v", expected, actual)
		}
	})

}

func TestIntegration(t *testing.T) {
	if !*integration {
		t.Skip("skipping integration tests")
	}

	r := NewDynamoRepository(nil)

	t.Run("success", func(t *testing.T) {
		var actual, expected error = nil, r.LoadTransaction()

		if actual != expected {
			t.Errorf("expected: %v, actual: %v", expected, actual)
		}
	})
}
