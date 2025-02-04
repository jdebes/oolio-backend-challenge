package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"jdebes/oolio-backend-challenge/repository"
)

func TestGetProduct(t *testing.T) {
	tests := []struct {
		name           string
		productId      string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Invalid Product ID",
			productId:      "invalid-id",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"description":"Invalid product ID"}`,
		},
		{
			name:           "Product Not Found",
			productId:      "9999",
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"description":"Product not found"}`,
		},
		{
			name:           "Success",
			productId:      "1",
			expectedStatus: http.StatusOK,
			expectedBody:   `{ "id": "1",  "name": "Waffle with Berries",  "category": "Waffle",  "price": 6.5}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			decimal.MarshalJSONWithoutQuotes = true
			mockProductDB := repository.NewProductDB()
			handler := BaseHandler{
				productDB: mockProductDB,
			}

			req, err := http.NewRequest("GET", "/api/product/"+tt.productId, nil)
			assert.NoError(t, err)

			rr := httptest.NewRecorder()
			router := mux.NewRouter()
			router.HandleFunc("/api/product/{productId}", handler.GetProduct)
			router.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			assert.JSONEq(t, tt.expectedBody, rr.Body.String())
		})
	}
}
