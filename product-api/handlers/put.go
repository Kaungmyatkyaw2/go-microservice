package handlers

import (
	"net/http"

	"github.com/Kaungmyatkyaw2/product-api/data"
)

// swagger:route PUT /products products updateProduct
// Update a products details
//
// responses:
//	201: noContentResponse
//  404: errorResponse
//  422: errorValidation

func (p *Products) UpdateProduct(w http.ResponseWriter, r *http.Request) {

	prod, ok := r.Context().Value(KeyProduct{}).(data.Product)

	if !ok {
		http.Error(w, "Failed to convert product", http.StatusBadRequest)
		return
	}

	err := data.UpdateProduct(prod.ID, &prod)

	if err == data.ErrProductNotFound {
		p.l.Println("[ERROR] product not found", err)

		w.WriteHeader(http.StatusNotFound)
		data.ToJSON(&GenericError{Message: "Product not found in database"}, w)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}
