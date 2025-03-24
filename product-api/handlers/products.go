package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	protos "github.com/Kaungmyatkyaw2/go-microservice/currency/protos/currency"
	"github.com/Kaungmyatkyaw2/go-microservice/product-api/data"

	"github.com/gorilla/mux"
)

type Products struct {
	l  *log.Logger
	v  *data.Validation
	cc protos.CurrencyClient
}

func NewProducts(l *log.Logger, v *data.Validation, cc protos.CurrencyClient) *Products {
	return &Products{l, v, cc}
}

var ErrInvalidProductPath = fmt.Errorf("Invalid Path, path should be /products/[id]")

type GenericError struct {
	Message string `json:"message"`
}

type ValidationError struct {
	Messages []string `json:"messages"`
}

func getProductID(r *http.Request) int {
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		panic(err)
	}

	return id
}
