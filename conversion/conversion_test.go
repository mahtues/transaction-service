package conversion

import (
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/pkg/errors"
)

type mockHttpClient func(request *http.Request) (*http.Response, error)

func (m mockHttpClient) Do(request *http.Request) (*http.Response, error) {
	return m(request)
}

func TestGetRateUnit(t *testing.T) {
	date, _ := time.Parse(time.RFC3339, "2022-07-31T20:04:59Z")

	const responseOneResult = `{"data":[{"record_date":"2022-06-30","exchange_rate":"5.182","country":"Brazil"}],"meta":{"count":1,"labels":{"record_date":"Record Date","exchange_rate":"Exchange Rate","country":"Country"},"dataTypes":{"record_date":"DATE","exchange_rate":"NUMBER","country":"STRING"},"dataFormats":{"record_date":"YYYY-MM-DD","exchange_rate":"10.2","country":"String"},"total-count":2,"total-pages":2},"links":{"self":"&page%5Bnumber%5D=1&page%5Bsize%5D=1","first":"&page%5Bnumber%5D=1&page%5Bsize%5D=1","prev":null,"next":"&page%5Bnumber%5D=2&page%5Bsize%5D=1","last":"&page%5Bnumber%5D=2&page%5Bsize%5D=1"}}`
	const responseNoResults = `{"data":[],"meta":{"count":0,"labels":{"record_date":"Record Date","exchange_rate":"Exchange Rate","country":"Country"},"dataTypes":{"record_date":"DATE","exchange_rate":"NUMBER","country":"STRING"},"dataFormats":{"record_date":"YYYY-MM-DD","exchange_rate":"10.2","country":"String"},"total-count":0,"total-pages":0},"links":{"self":"&page%5Bnumber%5D=1&page%5Bsize%5D=1","first":"&page%5Bnumber%5D=1&page%5Bsize%5D=1","prev":null,"next":"&page%5Bnumber%5D=2&page%5Bsize%5D=1","last":"&page%5Bnumber%5D=0&page%5Bsize%5D=1"}}`

	t.Run("connection error", func(t *testing.T) {
		var conv Service

		conv.Init(mockHttpClient(func(request *http.Request) (*http.Response, error) {
			return nil, errors.New("connection error")
		}))

		_, err := conv.GetRate("Brazil", date)
		if err == nil {
			t.Error("expected temporaraly unavailable")
		}
	})

	t.Run("error status code from remote api", func(t *testing.T) {
		var conv Service

		conv.Init(mockHttpClient(func(request *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusBadRequest,
				Body:       io.NopCloser(strings.NewReader("")),
			}, nil
		}))

		_, err := conv.GetRate("Brazil", date)
		if err == nil {
			t.Error("expected temporarily unavailable")
		}
	})

	t.Run("no results", func(t *testing.T) {
		var conv Service

		conv.Init(mockHttpClient(func(request *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(responseNoResults)),
			}, nil
		}))

		_, err := conv.GetRate("Brazil", date)
		if err == nil {
			t.Error("expected no results error")
		}
	})

	t.Run("success", func(t *testing.T) {
		var conv Service

		conv.Init(mockHttpClient(func(request *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(responseOneResult)),
			}, nil
		}))

		expected := "5.182"

		actual, err := conv.GetRate("Brazil", date)
		if err != nil {
			t.Error("expected no errors")
		}

		if expected != actual {
			t.Errorf("expected: %v, actual: %v", expected, actual)
		}
	})
}
