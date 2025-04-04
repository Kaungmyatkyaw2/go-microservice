package server

import (
	"context"
	"io"
	"time"

	"github.com/Kaungmyatkyaw2/go-microservice/currency/data"
	protos "github.com/Kaungmyatkyaw2/go-microservice/currency/protos/currency"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/go-hclog"
)

type Currency struct {
	rates         *data.ExchangeRates
	log           hclog.Logger
	subscriptions map[protos.Currency_SubscribeRatesServer][]*protos.RateRequest
	protos.UnimplementedCurrencyServer
}

func NewCurrency(r *data.ExchangeRates, l hclog.Logger) *Currency {

	c := &Currency{r, l, make(map[protos.Currency_SubscribeRatesServer][]*protos.RateRequest, 0), protos.UnimplementedCurrencyServer{}}
	go c.handleUpdates()
	return c
}

func (c *Currency) handleUpdates() {
	ru := c.rates.MonitorRates(5 * time.Second)
	for range ru {
		c.log.Info("Got updated rates")
		for k, v := range c.subscriptions {
			for _, rr := range v {
				r, err := c.rates.GetRate(rr.GetBase().String(), rr.GetDestination().String())
				if err != nil {
					c.log.Error("Unable to get updated rates", "base", rr.GetBase().String(), "destination", rr.GetDestination().String())
				}

				err = k.Send(&protos.StreamingRateResponse{Message: &protos.StreamingRateResponse_RateResponse{RateResponse: &protos.RateResponse{Base: rr.Base, Destination: rr.Destination, Rate: r}}})
				if err != nil {
					c.log.Error("Unable to send updated rates", "base", rr.GetBase().String(), "destination", rr.GetDestination().String())
				}
			}
		}
	}
}

func (c *Currency) GetRate(ctx context.Context, rr *protos.RateRequest) (*protos.RateResponse, error) {
	c.log.Info("Handle GetRate", "base", rr.GetBase(), "destination", rr.GetDestination())
	if rr.Base == rr.Destination {
		err := status.Newf(codes.InvalidArgument,
			"Base currency can not be the same as the destination currency",
		)

		err, wde := err.WithDetails(rr)

		if wde != nil {
			return nil, wde
		}
		return nil, err.Err()
	}
	rate, err := c.rates.GetRate(rr.GetBase().String(), rr.GetDestination().String())

	if err != nil {
		return nil, err
	}

	return &protos.RateResponse{Base: rr.Base, Destination: rr.Destination, Rate: rate}, nil
}

func (c *Currency) SubscribeRates(src protos.Currency_SubscribeRatesServer) error {
	for {
		rr, err := src.Recv()
		if err == io.EOF {
			c.log.Info("Client has closed connection")
			break
		}
		if err != nil {
			c.log.Error("Unable to read from client", "error", err)
			return err
		}

		c.log.Info("Handle client request", "request", rr)
		rrs, ok := c.subscriptions[src]

		if !ok {
			rrs = []*protos.RateRequest{}
		}

		// check that subscription does not exist

		for _, r := range rrs {
			// if we already have subscribe to this currency return an error
			if r.Base == rr.Base && r.Destination == rr.Destination {
				c.log.Error("Subscription already active", "base", rr.Base.String(), "dest", rr.Destination.String())

				grpcError := status.New(codes.InvalidArgument, "Subscription already active for rate")
				grpcError, err = grpcError.WithDetails(rr)
				if err != nil {
					c.log.Error("Unable to add metadata to error message", "error", err)
					continue
				}

				// Can't return error as that will terminate the connection, instead must send an error which
				// can be handled by the client Recv stream.
				rrs := &protos.StreamingRateResponse_Error{Error: grpcError.Proto()}
				src.Send(&protos.StreamingRateResponse{Message: rrs})
			}
		}

		rrs = append(rrs, rr)
		c.subscriptions[src] = rrs

	}
	return nil

}
