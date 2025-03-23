package handlers

import (
	"net/http"

	"github.com/Kaungmyatkyaw2/product-api/data"
)

// swagger:route GET /products products listProducts
// Returns a list of products
// responses:
// 	200: productsResponse

func (p *Products) ListAll(w http.ResponseWriter, _ *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	lp := data.GetProducts()

	err := data.ToJSON(lp, w)

	if err != nil {
		http.Error(w, "Unable to marshal json", http.StatusInternalServerError)
		return
	}
}

// swagger:route GET /products/{id} products getProductByID
// Return a specific product by matched id
// responses:
//	200: productResponse
//	404: errorResponse

// GetProductById handles GET requests
func (p *Products) GetByID(w http.ResponseWriter, r *http.Request) {

	id := getProductID(r)

	prod, err := data.GetProductByID(id)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)

		data.ToJSON(prod, w)

		return
	}

	data.ToJSON(prod, w)

}
