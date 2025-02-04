package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"jdebes/oolio-backend-challenge/models"
	"jdebes/oolio-backend-challenge/util"
)

func (s *BaseHandler) PlaceOrder(w http.ResponseWriter, r *http.Request) {
	var req models.OrderRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		util.WriteError(w, http.StatusBadRequest, util.InvalidInput)
		return
	}

	err = s.validator.Struct(&req)
	if err != nil {
		util.WriteError(w, http.StatusUnprocessableEntity, util.ValidationException)
		return
	}

	if !s.promoDB.IsValidPromoCode(req.CouponCode) {
		util.WriteError(w, http.StatusUnprocessableEntity, util.ValidationException)
		return
	}

	products := make([]models.Product, 0, len(req.Items))
	for _, item := range req.Items {
		product, exists := s.productDB.Get(item.ProductID)
		if !exists {
			util.WriteError(w, http.StatusUnprocessableEntity, util.ValidationException)
			return
		}

		products = append(products, product)
	}

	order := models.Order{
		ID:       uuid.New().String(),
		Items:    req.Items,
		Products: products,
	}

	if err := json.NewEncoder(w).Encode(order); err != nil {
		util.WriteError(w, http.StatusInternalServerError, util.InternalServerError)
		log.Errorf("Failed to encode order to JSON: %v", err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}
