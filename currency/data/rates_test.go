package data

import (
	"fmt"
	"testing"

	"github.com/hashicorp/go-hclog"
)

func TestNewRates(t *testing.T) {

	tr, err := NewExchangeRates(hclog.Default())

	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("%#v", tr.rates)

}
