package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Kaungmyatkyaw2/go-microservice/product-api/data"
	"github.com/hashicorp/go-hclog"

	"github.com/gorilla/mux"
)

type Products struct {
	l          hclog.Logger
	v          *data.Validation
	productsDB *data.ProductsDB
}

func NewProducts(l hclog.Logger, v *data.Validation, pdb *data.ProductsDB) *Products {
	return &Products{l, v, pdb}
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
