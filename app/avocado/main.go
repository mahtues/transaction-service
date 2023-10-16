package main

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"

	"github.com/mahtues/transaction-service/app/avocado/server"
	"github.com/mahtues/transaction-service/conversion"
	"github.com/mahtues/transaction-service/support"
	"github.com/mahtues/transaction-service/transaction"
)

type Avocado struct {
	// resources
	awsResources support.AwsResources

	// repositories
	transactionRepository      transaction.DynamoRepository
	transactionInMemRepository transaction.InMemRepository

	// services
	transactionService transaction.Service
	conversionService  conversion.Service

	// todo: configs
	// todo: loggers
	// todo: monitoring
	// todo: warmup

	httpServer server.Server
}

func NewAvocado() *Avocado {
	avocado := &Avocado{}

	// resources
	// avocado.awsResources.Init()

	// repositories
	// avocado.transactionRepository.Init(&avocado.awsResources, "transactions")

	avocado.transactionInMemRepository.Init()

	// services
	avocado.transactionService.Init(
		&avocado.transactionInMemRepository,
		&avocado.conversionService,
	)

	avocado.conversionService.Init(
		http.DefaultClient,
	)

	avocado.httpServer.Init(
		&avocado.transactionService,
	)

	return avocado
}

func (a *Avocado) Start() error {
	if err := a.httpServer.Start(); err != nil {
		return errors.Wrap(err, "failed to start application")
	}

	return nil
}

func main() {
	fmt.Println("starting avocado application...")

	avocado := NewAvocado()

	if err := avocado.Start(); err != nil {
		fmt.Println("exiting avocado application due to error:", err)
	}
}
