package handlers

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func (s *BaseHandler) ListProducts(w http.ResponseWriter, r *http.Request) {
	products := s.db.List()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(products)
	if err != nil {
		log.Errorf("Failed to encode products to JSON: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
