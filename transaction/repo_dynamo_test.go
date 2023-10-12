package transaction

import (
	"flag"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

var integration = flag.Bool("integration", false, "only perform local tests")

func TestIntegrationDynamoRepo(t *testing.T) {
	if !*integration {
		t.Skip("skipping integration tests")
	}

	sess, err := session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Region:      aws.String("us-east-2"),
			Endpoint:    aws.String("http://localhost:4579"),
			Credentials: credentials.NewStaticCredentials("fake-access-key", "fake-secret-key", "fake-token"),
		},
		SharedConfigState: session.SharedConfigEnable,
	})

	if err != nil {
		t.Fatal("failed to create aws session")
	}

	dyna := dynamodb.New(sess)
	if dyna == nil {
		t.Fatal("dynamo client is nil")
	}

	tableName := RandStringRunes(10) + "-transactions"

	if err = CreateTable(dyna, tableName); err != nil {
		t.Fatal(err)
	}
	//defer DeleteTable(dyna, tableName)

	repo := NewDynamoRepository(dyna, tableName)

	t.Run("save transaction", func(t *testing.T) {
		transaction := Transaction{
			Id:          "tr-" + RandStringRunes(20),
			Description: "short description",
		}

		var actual, expected error = repo.SaveTransaction(transaction), nil

		if actual != expected {
			t.Errorf("expected: %v, actual: %v", expected, actual)
		}
	})

	t.Run("load transaction", func(t *testing.T) {
		expected := Transaction{
			Id:          "tr-" + RandStringRunes(20),
			Description: "short description",
			Date:        time.Now().Truncate(time.Second),
			AmountUs:    "123.30",
		}

		{
			err := repo.SaveTransaction(expected)
			if err != nil {
				t.Errorf("err: %v", err)
			}
		}

		{
			actual, err := repo.LoadTransaction(expected.Id)
			if err != nil {
				t.Errorf("expected: %v, err: %v", expected, err)
			}
			if actual != expected {
				t.Errorf("expected: %v, actual: %v", expected, actual)
			}
		}
	})
}
