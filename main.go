package main

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
	"jdebes/oolio-backend-challenge/data"
	"jdebes/oolio-backend-challenge/handlers"
	"jdebes/oolio-backend-challenge/middleware"
)

const (
	dataArg   = "data"
	serverArg = "server"
)

func main() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetLevel(log.DebugLevel)

	if len(os.Args) < 2 {
		log.Panicf("Requires argument '%s' or '%s'", dataArg, serverArg)
		os.Exit(1)
	}

	switch cmd := os.Args[1]; cmd {
	case dataArg:
		err := data.WriteValidPromos()
		if err != nil {
			log.Panic("Failed write valid promos", err)
		}
	case serverArg:
		// This decimal values are returned as a number in JSON. This brings functionality in line with api spec.
		// However, this is not the best idea, as it risks clients suffering from precision loss.
		decimal.MarshalJSONWithoutQuotes = true

		base := handlers.NewBaseHandler()

		r := mux.NewRouter().StrictSlash(true)

		public := r.PathPrefix("/api").Subrouter()
		public.HandleFunc("/product/{productId}", base.GetProduct).Methods(http.MethodGet)
		public.HandleFunc("/product", base.ListProducts).Methods(http.MethodGet)

		auth := r.PathPrefix("/api").Subrouter()
		auth.Use(middleware.AuthMiddleware)
		auth.HandleFunc("/order", base.PlaceOrder).Methods(http.MethodPost)

		log.Info("Server starting on port 8080...")
		log.Fatal(http.ListenAndServe(":8080", r))
	default:
		log.Errorf("Invalid argument, provide '%s' or '%s'", dataArg, serverArg)
		os.Exit(1)
	}
}
