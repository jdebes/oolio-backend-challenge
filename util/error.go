package util

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
)

const (
	InvalidInput        = "Invalid input"
	ValidationException = "Validation exception"
	InternalServerError = "Internal server error"
	ProductNotFound     = "Product not found"
	InvalidProductID    = "Invalid product ID"
)

type ErrorResponse struct {
	Description string `json:"description"`
}

func WriteError(w http.ResponseWriter, statusCode int, description string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	errResponse := ErrorResponse{Description: description}
	if err := json.NewEncoder(w).Encode(errResponse); err != nil {
		log.Errorf("Failed to encode error response: %v", err)
	}
}
