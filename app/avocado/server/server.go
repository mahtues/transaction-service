package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/mahtues/form"
	"github.com/mahtues/transaction-service/apperrors"
	"github.com/mahtues/transaction-service/support"
	"github.com/mahtues/transaction-service/transaction"
)

type Server struct {
	transactionService *transaction.Service

	handler http.Handler
}

func (s *Server) Init(transactionService *transaction.Service) {
	s.transactionService = transactionService

	m := http.NewServeMux()

	m.HandleFunc("/transaction", MustMethodFunc(http.MethodPost, s.createTransaction))
	m.HandleFunc("/transaction/", MustMethodFunc(http.MethodGet, s.getTransaction))
	m.HandleFunc("/heartbeat", s.heartbeat)

	s.handler = m
}

func (s *Server) heartbeat(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "alive")
}

func (s *Server) getTransaction(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/transaction/"):]
	country := r.FormValue("country")

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
		Id:      id,
		Country: country,
	}

	response := transaction.GetResponse{}

	if response, err = s.transactionService.GetTransaction(request); err != nil {
		apperror := apperrors.Cause(err)

		json.NewEncoder(w).Encode(respError{
			Error: apperror.Error(),
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
}

func (s *Server) createTransaction(w http.ResponseWriter, r *http.Request) {
	frm := struct {
		Description *string       `form:"description"`
		Date        *support.Time `form:"date"`
		AmountUs    *string       `form:"amountUs"`
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

	invalid := []string{}

	if frm.Description == nil {
		invalid = append(invalid, "description missing")
	} else if len(*frm.Description) == 0 || len(*frm.Description) > 50 {
		invalid = append(invalid, "description must contain between 1 and 50 characters")
	}

	if frm.Date == nil {
		invalid = append(invalid, "date missing")
	}

	if frm.AmountUs == nil || len(*frm.AmountUs) == 0 {
		invalid = append(invalid, "amountUs missing")
	}

	if len(invalid) > 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(respError{
			Error: fmt.Sprintf("invalid form: %s", strings.Join(invalid, ", ")),
		})
		return
	}

	request := transaction.CreateRequest{
		Description: *frm.Description,
		Date:        time.Time(*frm.Date),
		AmountUs:    *frm.AmountUs,
	}

	response := transaction.CreateResponse{}

	if response, err = s.transactionService.CreateTransaction(request); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		apperr := apperrors.Cause(err)
		json.NewEncoder(w).Encode(respError{
			Error: apperr.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(respOk{
		Id: response.Id,
	})
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.handler.ServeHTTP(w, r)
}

func (s *Server) Start() error {
	if err := http.ListenAndServe(":8000", s); err != nil {
		return errors.Wrap(err, "failed to start server")
	}

	return nil
}
