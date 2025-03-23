package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/Kaungmyatkyaw2/go-microservice/product-api/data"
	"github.com/Kaungmyatkyaw2/go-microservice/product-api/handlers"

	"github.com/go-openapi/runtime/middleware"
	gohandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"google.golang.org/grpc"
)

func main() {

	l := log.New(os.Stdout, "product-api: ", log.LstdFlags)
	v := data.NewValidation()

	conn, err := grpc.Dial("localhost:9092")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// cc := protos.NewCurrenctyClient(conn)

	ph := handlers.NewProducts(l, v)

	sm := mux.NewRouter()
	getRouter := sm.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/products", ph.ListAll)

	putRouter := sm.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/products", ph.UpdateProduct)
	putRouter.Use(ph.MiddlewareProductValidation)

	addRouter := sm.Methods(http.MethodPost).Subrouter()
	addRouter.HandleFunc("/products", ph.Create)
	addRouter.Use(ph.MiddlewareProductValidation)

	opts := middleware.RedocOpts{SpecURL: "./swagger.yaml"}
	sh := middleware.Redoc(opts, nil)
	getRouter.Handle("/docs", sh)

	getRouter.Handle("/swagger.yaml", http.FileServer(http.Dir("./")))

	ch := gohandlers.CORS(gohandlers.AllowedOrigins([]string{"http://localhost:3000"}))

	s := &http.Server{
		Addr:         ":9090",
		Handler:      ch(sm),
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	go func() {
		err := s.ListenAndServe()

		if err != nil {
			l.Fatal(err)
		}
	}()

	sigChan := make(chan os.Signal)

	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <-sigChan

	l.Println("Received terminate, graceful shutdown", sig)

	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)

	s.Shutdown(tc)

}
