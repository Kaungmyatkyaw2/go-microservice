package handlers

import (
	"net/http"

	"github.com/Kaungmyatkyaw2/go-microservice/product-api/data"
)

// swagger:route DELETE /products/{id} products deleteProduct
// Deleting a product
//
// responses:
//	201: noContentResponse

func (p *Products) Delete(rw http.ResponseWriter, r *http.Request) {
	id := getProductID(r)

	p.l.Debug("Deleting record id", id)

	err := p.productsDB.DeleteProduct(id)
	if err == data.ErrProductNotFound {
		p.l.Error("Deleting record id does not exist")

		rw.WriteHeader(http.StatusNotFound)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)
		return
	}

	if err != nil {
		p.l.Error("Deleting record", "error", err)

		rw.WriteHeader(http.StatusInternalServerError)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}
