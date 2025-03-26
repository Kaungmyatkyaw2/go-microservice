package handlers

import (
	"net/http"

	"github.com/Kaungmyatkyaw2/go-microservice/product-api/data"
)

// swagger:route GET /products products listProducts
// Returns a list of products
// responses:
// 	200: productsResponse

func (p *Products) ListAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	cur := r.URL.Query().Get("currency")

	lp, err := p.productsDB.GetProducts(cur)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		data.ToJSON(&GenericError{Message: err.Error()}, w)
		return
	}

	err = data.ToJSON(lp, w)
	if err != nil {
		p.l.Error("Unable to serializing product", "error", err)
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
	cur := r.URL.Query().Get("currency")

	p.l.Debug("Currency Query", cur)

	prod, err := p.productsDB.GetProductByID(id, cur)

	switch err {
	case nil:

	case data.ErrProductNotFound:
		p.l.Error("Unable to fetch product", "error", err)

		w.WriteHeader(http.StatusNotFound)
		data.ToJSON(&GenericError{Message: err.Error()}, w)
		return
	default:
		p.l.Error("Unable to fetching product", "error", err)

		w.WriteHeader(http.StatusInternalServerError)
		data.ToJSON(&GenericError{Message: err.Error()}, w)
		return
	}

	err = data.ToJSON(prod, w)
	if err != nil {
		// we should never be here but log the error just incase
		p.l.Error("Unable to serializing product", err)
	}
}
