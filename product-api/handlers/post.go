package handlers

import (
	"net/http"

	"github.com/Kaungmyatkyaw2/go-microservice/product-api/data"
)

// swagger:route POST /products products createProduct
// Create a new product
//
// responses:
//
//		200: productResponse
//	 422: errorValidation
//	 501: errorResponse
func (p *Products) Create(w http.ResponseWriter, r *http.Request) {
	prod, ok := r.Context().Value(KeyProduct{}).(data.Product)

	if !ok {
		http.Error(w, "Failed to convert product", http.StatusBadRequest)
		return
	}

	p.productsDB.AddProduct(&prod)

}
