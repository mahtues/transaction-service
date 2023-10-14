package main

import (
	"flag"
	"testing"
	"time"
)

var integration = flag.Bool("integration", false, "only perform local tests")

func TestMain(t *testing.T) {
	avocado := NewAvocado()

	avocado.conversionService.GetRate("", time.Time{})
}
