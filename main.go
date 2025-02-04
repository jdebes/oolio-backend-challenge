package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
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
		base := handlers.NewBaseHandler()

		r := mux.NewRouter()

		r.Use(middleware.AuthMiddleware)

		r.HandleFunc("/product/{productId}", base.GetProduct).Methods(http.MethodGet)
		r.HandleFunc("/product", base.ListProducts).Methods(http.MethodGet)
		r.HandleFunc("/order", base.PlaceOrder).Methods(http.MethodPost)

		fmt.Println("Server starting on port 8080...")
		log.Fatal(http.ListenAndServe(":8080", r))
	default:
		log.Errorf("Invalid argument, provide '%s' or '%s'", dataArg, serverArg)
		os.Exit(1)
	}
}
