package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"jdebes/oolio-backend-challenge/models"
)

func (s *BaseHandler) PlaceOrder(w http.ResponseWriter, r *http.Request) {
	var req models.OrderRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error parsing request: %s", err), http.StatusBadRequest)
		return
	}

	err = s.validator.Struct(req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Order Validation failed: %s", err), http.StatusBadRequest)
		return
	}

	// TODO validate coupon code function for advanced validation

	var products []models.Product
	for _, item := range req.Items {
		product, exists := s.db.Get(item.ProductID)
		if !exists {
			http.Error(w, fmt.Sprintf("Product not found: %s", item.ProductID), http.StatusBadRequest)
			return
		}

		products = append(products, product)
	}

	order := models.Order{
		ID:       uuid.New().String(),
		Items:    req.Items,
		Products: products,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(order); err != nil {
		http.Error(w, fmt.Sprintf("Error sending response: %s", err), http.StatusInternalServerError)
	}
}
