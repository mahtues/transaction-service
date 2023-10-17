package conversion

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/mahtues/transaction-service/apperrors"
	"github.com/mahtues/transaction-service/support"
	"github.com/pkg/errors"
)

const (
	requestUrl = "https://api.fiscaldata.treasury.gov/services/api/fiscal_service/v1/accounting/od/rates_of_exchange"
	filterFmt  = "country:eq:%s,record_date:lte:%s,record_date:gt:%s" // country, record date, six months before
)

type Service struct {
	httpClient httpClient
}

type httpClient interface {
	Do(request *http.Request) (*http.Response, error)
}

func (s *Service) Init(httpClient httpClient) {
	s.httpClient = httpClient
}

func (s *Service) GetRate(country string, date time.Time) (string, error) {
	req, err := http.NewRequest(http.MethodGet, requestUrl, nil)
	if err != nil {
		return "", errors.Wrap(err, "unexpected error")
	}

	q := url.Values{}
	q.Set("fields", "record_date,exchange_rate,country")
	q.Set("filter", fmt.Sprintf(filterFmt, country, date.Format(time.DateOnly), date.AddDate(0, -6, 0).Format(time.DateOnly)))
	q.Set("sort", "-record_date")
	q.Set("page[size]", "1")

	req.URL.RawQuery = q.Encode()

	res, err := s.httpClient.Do(req)
	if err != nil {
		// return "", errors.Wrap(err, "failed to do request")
		return "", apperrors.Wrapf(errors.New(err.Error()), "service temporarily unavailable")
	}

	if res.StatusCode != http.StatusOK {
		// return "", errors.New(fmt.Sprint("response status code: ", res.Status))
		return "", apperrors.Wrapf(errors.New(res.Status), "service temporarily unavailable")
	}

	type rateEntry struct {
		RecordDate   support.Date `json:"record_date"`
		ExchangeRate string       `json:"exchange_rate"`
		Country      string       `json:"country"`
	}

	type apiResponse struct {
		Data []rateEntry `json:"data"`
	}

	var resBody apiResponse

	err = json.NewDecoder(res.Body).Decode(&resBody)
	if err != nil {
		// return "", errors.Wrap(err, "failed to read response body")
		return "", apperrors.Wrapf(errors.New(err.Error()), "service temporarily unavailable")
	}

	if len(resBody.Data) == 0 {
		// return "", errors.New("no rates found for date range")
		return "", apperrors.Wrapf(errors.New(""), "no rates found for date range for specified country")
	}

	return resBody.Data[0].ExchangeRate, nil
}

func Convert(value string, rate string) (string, error) {
	v, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return "", err
	}

	r, err := strconv.ParseFloat(rate, 64)
	if err != nil {
		return "", err
	}

	c := v * r
	converted := fmt.Sprintf("%.2f", c)

	return converted, nil
}
