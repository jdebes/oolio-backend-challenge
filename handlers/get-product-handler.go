package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"jdebes/oolio-backend-challenge/util"
)

func (s *BaseHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productId := vars["productId"]

	id, err := strconv.Atoi(productId)
	if err != nil || id <= 0 {
		util.WriteError(w, http.StatusBadRequest, util.InvalidProductID)
		return
	}

	product, exists := s.productDB.Get(productId)
	if !exists {
		util.WriteError(w, http.StatusNotFound, util.ProductNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(product)
	if err != nil {
		util.WriteError(w, http.StatusInternalServerError, util.InternalServerError)
		log.Errorf("Failed to encode product to JSON: %v", err)
	}
}
