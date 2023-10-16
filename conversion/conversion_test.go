package conversion

import (
	"net/http"
	"testing"
	"time"
)

func TestConversion(t *testing.T) {
	var conv Service

	conv.Init(http.DefaultClient)

	date, _ := time.Parse(time.DateOnly, "2022-08-30")
	rate, err := conv.GetRate("Brazil", date)
	t.Log(rate, err)
}
