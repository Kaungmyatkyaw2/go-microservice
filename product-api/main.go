package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/Kaungmyatkyaw2/go-microservice/product-api/data"
	"github.com/Kaungmyatkyaw2/go-microservice/product-api/handlers"
	"github.com/hashicorp/go-hclog"

	"github.com/go-openapi/runtime/middleware"
	gohandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"google.golang.org/grpc"

	protos "github.com/Kaungmyatkyaw2/go-microservice/currency/protos/currency"
)

func main() {

	l := hclog.Default()
	v := data.NewValidation()

	conn, err := grpc.Dial("localhost:9092", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	cc := protos.NewCurrencyClient(conn)
	pdb := data.NewProductsDB(cc, l)
	ph := handlers.NewProducts(l, v, pdb)

	sm := mux.NewRouter()
	getRouter := sm.Methods(http.MethodGet).Subrouter()

	getRouter.HandleFunc("/products", ph.ListAll).Queries("currency", "{[A-Z]{3}}")
	getRouter.HandleFunc("/products", ph.ListAll)

	getRouter.HandleFunc("/products/{id:[0-9]+}", ph.GetByID).Queries("currency", "{[A-Z]{3}}")
	getRouter.HandleFunc("/products/{id:[0-9]+}", ph.GetByID)

	putRouter := sm.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/products", ph.UpdateProduct)
	putRouter.Use(ph.MiddlewareProductValidation)

	addRouter := sm.Methods(http.MethodPost).Subrouter()
	addRouter.HandleFunc("/products", ph.Create)
	addRouter.Use(ph.MiddlewareProductValidation)

	deleteR := sm.Methods(http.MethodDelete).Subrouter()
	deleteR.HandleFunc("/products/{id:[0-9]+}", ph.Delete)

	opts := middleware.RedocOpts{SpecURL: "./swagger.yaml"}
	sh := middleware.Redoc(opts, nil)
	getRouter.Handle("/docs", sh)

	getRouter.Handle("/swagger.yaml", http.FileServer(http.Dir("./")))

	ch := gohandlers.CORS(gohandlers.AllowedOrigins([]string{"http://localhost:3000"}))

	s := &http.Server{
		Addr:         ":9090",
		Handler:      ch(sm),
		IdleTimeout:  120 * time.Second,
		ErrorLog:     l.StandardLogger(&hclog.StandardLoggerOptions{}),
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	go func() {
		l.Info("Starting server on port 9090")

		err := s.ListenAndServe()
		if err != nil {
			l.Error("Error starting server", "error", err)
			os.Exit(1)
		}
	}()
	sigChan := make(chan os.Signal)

	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <-sigChan

	l.Info("Received terminate, graceful shutdown", sig)

	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)

	s.Shutdown(tc)

}
