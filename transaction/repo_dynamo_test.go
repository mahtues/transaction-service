package transaction

import (
	"flag"
	"testing"
	"time"

	"github.com/mahtues/transaction-service/random"
	"github.com/mahtues/transaction-service/support"
)

var integration = flag.Bool("integration", false, "only perform local tests")

func TestIntegrationDynamoRepository(t *testing.T) {
	if !*integration {
		t.Skip("skipping integration tests")
	}

	var err error

	var awsResources support.AwsResources

	awsResources.Init()

	tableName := random.String(10) + "-transactions"

	if err = CreateTable(awsResources.DynamoDbClient, tableName); err != nil {
		t.Fatal(err)
	}
	defer DeleteTable(awsResources.DynamoDbClient, tableName)

	repository := NewDynamoRepository(&awsResources, tableName)

	t.Run("save transaction", func(t *testing.T) {
		transaction := Transaction{
			Id:          "tr-" + random.String(20),
			Description: "short description",
		}

		var actual, expected error = repository.SaveTransaction(transaction), nil

		if actual != expected {
			t.Errorf("expected: %v, actual: %v", expected, actual)
		}
	})

	t.Run("load transaction", func(t *testing.T) {
		expected := Transaction{
			Id:          "tr-" + random.String(20),
			Description: "short description",
			Date:        time.Now().Truncate(time.Second),
			AmountUs:    "123.30",
		}

		err := repository.SaveTransaction(expected)
		if err != nil {
			t.Errorf("err: %v", err)
		}

		actual, err := repository.LoadTransaction(expected.Id)
		if err != nil {
			t.Errorf("expected: %v, err: %v", nil, err)
		}
		if !actual.Equal(expected) {
			t.Errorf("   expected: %v, actual: %v", expected, actual)
		}
	})
}
