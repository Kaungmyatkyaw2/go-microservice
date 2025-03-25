package main

import (
	"net"
	"os"

	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/Kaungmyatkyaw2/go-microservice/currency/data"
	protos "github.com/Kaungmyatkyaw2/go-microservice/currency/protos/currency"
	"github.com/Kaungmyatkyaw2/go-microservice/currency/server"
)

func main() {

	log := hclog.Default()

	rates, err := data.NewExchangeRates(log)

	if err != nil {
		log.Error("Unable to generate rates", "error", err)
		os.Exit(1)
	}

	gs := grpc.NewServer()

	cs := server.NewCurrency(rates, log)

	protos.RegisterCurrencyServer(gs, cs)

	reflection.Register(gs)

	l, err := net.Listen("tcp", ":9092")
	if err != nil {
		log.Error("Unable to listen", "error", err)
		os.Exit(1)
	}

	gs.Serve(l)

}
