package server

import (
	"context"

	"github.com/Kaungmyatkyaw2/go-microservice/currency/data"
	protos "github.com/Kaungmyatkyaw2/go-microservice/currency/protos/currency"

	"github.com/hashicorp/go-hclog"
)

type Currency struct {
	rates *data.ExchangeRates
	log   hclog.Logger
	protos.UnimplementedCurrencyServer
}

func NewCurrency(r *data.ExchangeRates, l hclog.Logger) *Currency {
	return &Currency{r, l, protos.UnimplementedCurrencyServer{}}
}

func (c *Currency) GetRate(ctx context.Context, rr *protos.RateRequest) (*protos.RateResponse, error) {
	c.log.Info("Handle GetRate", "base", rr.GetBase(), "destination", rr.GetDestination())

	rate, err := c.rates.GetRate(rr.GetBase().String(), rr.GetDestination().String())

	if err != nil {
		return nil, err
	}

	return &protos.RateResponse{Rate: rate}, nil
}
