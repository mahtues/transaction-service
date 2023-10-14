package main

import (
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"

	"github.com/mahtues/transaction-service/conversion"
	"github.com/mahtues/transaction-service/support"
	"github.com/mahtues/transaction-service/transaction"
)

type Avocado struct {
	// resources
	awsSession        *session.Session
	awsDynamoDbClient *dynamodb.DynamoDB

	// repositories
	transactionRepository transaction.DynamoRepository

	// services
	transactionService transaction.Service
	conversionService  conversion.Service

	// todo: configs
	// todo: loggers
	// todo: monitoring
	// todo: heartbeat
}

func NewAvocado() *Avocado {
	avocado := &Avocado{}

	var err error

	// resources
	if avocado.awsSession, err = support.AwsSessionLocalhost(4579); err != nil {
		panic("message")
	}

	if avocado.awsDynamoDbClient, err = support.AwsDynamoDbClient(avocado.awsSession); err != nil {
		panic("message")
	}

	// repositories
	avocado.transactionRepository = *transaction.NewDynamoRepository(
		avocado.awsDynamoDbClient,
		"prod-transaction",
	)

	// services
	avocado.transactionService = *transaction.NewService(
		&avocado.transactionRepository,
		&avocado.conversionService,
	)

	return avocado
}

func main() {
	fmt.Println("starting avocado application...")

	mux := http.NewServeMux()

	if err := http.ListenAndServe(":8000", mux); err != nil {
		fmt.Println("failed to start server")
	}
}
