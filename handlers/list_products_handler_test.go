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

func TestListProducts(t *testing.T) {
	tests := []struct {
		name           string
		mockProductDB  *repository.ProductDB
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Empty Result",
			mockProductDB:  repository.NewEmptyProductDB(),
			expectedStatus: http.StatusOK,
			expectedBody:   `[]`,
		},
		{
			name:           "All Products",
			mockProductDB:  repository.NewProductDB(),
			expectedStatus: http.StatusOK,
			expectedBody: `[
				{"id":"1","name":"Waffle with Berries","category":"Waffle","price":6.5},
				{"id":"2","name":"Vanilla Bean Crème Brûlée","category":"Crème Brûlée","price":7},
				{"id":"3","name":"Macaron Mix of Five","category":"Macaron","price":8},
				{"id":"4","name":"Classic Tiramisu","category":"Tiramisu","price":5.5},
				{"id":"5","name":"Pistachio Baklava","category":"Baklava","price":4}
			]`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			decimal.MarshalJSONWithoutQuotes = true

			handler := BaseHandler{
				productDB: tt.mockProductDB,
			}

			req, err := http.NewRequest("GET", "/api/product", nil)
			assert.NoError(t, err)

			rr := httptest.NewRecorder()
			router := mux.NewRouter()
			router.HandleFunc("/api/product", handler.ListProducts)
			router.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			assert.JSONEq(t, tt.expectedBody, rr.Body.String())
		})
	}
}
