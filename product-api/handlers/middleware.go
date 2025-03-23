package handlers

import (
	"context"
	"net/http"

	"github.com/Kaungmyatkyaw2/go-microservice/product-api/data"
)

type KeyProduct struct{}

func (p *Products) MiddlewareProductValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		prod := data.Product{}

		err := data.FromJSON(prod, r.Body)

		if err != nil {
			p.l.Println("[ERROR] desrializing product", err)
			w.WriteHeader(http.StatusBadRequest)
			data.ToJSON(&GenericError{Message: err.Error()}, w)
			return
		}

		errs := p.v.Validate(prod)
		if len(errs) != 0 {

			p.l.Println("[ERROR] validating product", errs)

			w.WriteHeader(http.StatusUnprocessableEntity)
			data.ToJSON(&ValidationError{Messages: errs.Errors()}, w)
		}

		ctx := context.WithValue(r.Context(), KeyProduct{}, prod)
		req := r.WithContext(ctx)

		next.ServeHTTP(w, req)
	})
}
