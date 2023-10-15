package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/pkg/errors"

	"github.com/mahtues/form"

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
	// todo: heartbeat
	// todo: warmup

	httpServer server
}

func NewAvocado() *Avocado {
	avocado := &Avocado{}

	// resources
	// avocado.awsResources.Init()

	// repositories
	// avocado.transactionRepository.Init(&avocado.awsResources, "transactions")

	avocado.transactionInMemRepository.Init()

	// services
	avocado.transactionService = *transaction.NewService(
		&avocado.transactionInMemRepository,
		&avocado.conversionService,
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

type server struct {
	transactionService *transaction.Service

	handler http.Handler
}

func (s *server) Init(transactionService *transaction.Service) {
	s.transactionService = transactionService

	m := http.NewServeMux()

	m.HandleFunc("/heartbeat", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "alive")
	})

	m.HandleFunc("/transaction", func(w http.ResponseWriter, r *http.Request) {
		frm := struct {
			Description string `form:"description"`
			Date        Date   `form:"date"`
			AmountUs    string `form:"amountUs"`
		}{}

		type respOk struct {
			Id string `json:"id,omitempty"`
		}

		type respError struct {
			Error string `json:"error,omitempty"`
		}

		var err error

		if err = form.Unmarshal(r, &frm); err != nil {
			json.NewEncoder(w).Encode(respError{
				Error: err.Error(),
			})
			return
		}

		request := transaction.CreateRequest{
			Description: frm.Description,
			Date:        time.Time(frm.Date),
			AmountUs:    frm.AmountUs,
		}

		response := transaction.CreateResponse{}

		if response, err = s.transactionService.CreateTransaction(request); err != nil {
			json.NewEncoder(w).Encode(respError{
				Error: err.Error(),
			})
			return
		}

		json.NewEncoder(w).Encode(respOk{
			Id: response.Id,
		})
	})

	m.HandleFunc("/transaction/", func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Path[len("/transaction/"):]
		currency := r.FormValue("curr")

		type respOk struct {
			Id              string    `json:"id"`
			Description     string    `json:"description"`
			Date            time.Time `json:"date"`
			AmountUs        string    `json:"amountUs"`
			ConversionRate  string    `json:"conversionRate"`
			AmountConverted string    `json:"amountConverted"`
		}

		type respError struct {
			Error string `json:"error,omitempty"`
		}

		var err error

		request := transaction.GetRequest{
			Id:       id,
			Currency: currency,
		}

		response := transaction.GetResponse{}

		if response, err = s.transactionService.GetTransaction(request); err != nil {
			json.NewEncoder(w).Encode(respError{
				Error: err.Error(),
			})
			return
		}

		json.NewEncoder(w).Encode(respOk{
			Id:              response.Id,
			Description:     response.Description,
			Date:            response.Date,
			AmountUs:        response.AmountUs,
			ConversionRate:  response.Rate,
			AmountConverted: response.AmountConverted,
		})
	})

	s.handler = m
}

type Date time.Time

func (d *Date) UnmarshalFormField(s string) error {
	t, err := time.Parse(time.RFC3339, s)

	if err != nil {
		return err
	}

	*d = Date(t)

	return nil
}

func (s *server) Start() error {
	if err := http.ListenAndServe(":8000", s.handler); err != nil {
		return errors.Wrap(err, "failed to start server")
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
