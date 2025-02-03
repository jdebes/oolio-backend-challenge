package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"jdebes/oolio-backend-challenge/handlers"
	"jdebes/oolio-backend-challenge/middleware"
)

func main() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetLevel(log.DebugLevel)

	base := handlers.NewBaseHandler()

	r := mux.NewRouter()

	r.Use(middleware.AuthMiddleware)

	r.HandleFunc("/product/{productId}", base.GetProduct).Methods(http.MethodGet)
	r.HandleFunc("/product", base.ListProducts).Methods(http.MethodGet)
	r.HandleFunc("/order", base.PlaceOrder).Methods(http.MethodPost)

	fmt.Println("Server starting on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}
