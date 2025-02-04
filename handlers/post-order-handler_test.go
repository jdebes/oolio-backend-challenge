package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"jdebes/oolio-backend-challenge/models"
	"jdebes/oolio-backend-challenge/repository"
)

func TestPlaceOrder(t *testing.T) {
	tests := []struct {
		name           string
		orderRequest   models.OrderRequest
		setupHandler   func() *BaseHandler
		expectedStatus int
		expectedBody   string
	}{
		{
			name:         "Validation Error",
			orderRequest: models.OrderRequest{},
			setupHandler: func() *BaseHandler {
				mockPromoDB := repository.NewEmptyPromoDB()
				productDB := repository.NewEmptyProductDB()

				return &BaseHandler{
					productDB: productDB,
					promoDB:   mockPromoDB,
					validator: validator.New(),
				}
			},
			expectedStatus: http.StatusUnprocessableEntity,
			expectedBody:   `{"description":"Validation exception"}`,
		},
		{
			name: "Validation Error - Missing Coupon Code",
			orderRequest: models.OrderRequest{
				Items: []models.LineItem{
					{ProductID: "1", Quantity: 1},
				},
			},
			setupHandler: func() *BaseHandler {
				mockPromoDB := repository.NewEmptyPromoDB()
				productDB := repository.NewEmptyProductDB()

				return &BaseHandler{
					productDB: productDB,
					promoDB:   mockPromoDB,
					validator: validator.New(),
				}
			},
			expectedStatus: http.StatusUnprocessableEntity,
			expectedBody:   `{"description":"Validation exception"}`,
		},
		{
			name: "Invalid Promo Code",
			orderRequest: models.OrderRequest{
				CouponCode: "INVALIDCODE",
				Items: []models.LineItem{
					{ProductID: "1", Quantity: 1},
				},
			},
			expectedStatus: http.StatusUnprocessableEntity,
			expectedBody:   `{"description":"Validation exception"}`,
			setupHandler: func() *BaseHandler {
				mockPromoDB := repository.NewEmptyPromoDB()
				productDB := repository.NewEmptyProductDB()

				return &BaseHandler{
					productDB: productDB,
					promoDB:   mockPromoDB,
					validator: validator.New(),
				}
			},
		},
		{
			name: "Product Not Found",
			orderRequest: models.OrderRequest{
				CouponCode: "VALIDCODE",
				Items: []models.LineItem{
					{ProductID: "9999", Quantity: 1}, // Non-existent product
				},
			},
			expectedStatus: http.StatusUnprocessableEntity,
			expectedBody:   `{"description":"Validation exception"}`,
			setupHandler: func() *BaseHandler {
				mockPromoDB := repository.NewEmptyPromoDB()
				productDB := repository.NewEmptyProductDB()

				return &BaseHandler{
					productDB: productDB,
					promoDB:   mockPromoDB,
					validator: validator.New(),
				}
			},
		},
		{
			name: "Order Created",
			orderRequest: models.OrderRequest{
				CouponCode: "VALIDCODE",
				Items: []models.LineItem{
					{ProductID: "1", Quantity: 1},
					{ProductID: "2", Quantity: 2},
				},
			},
			setupHandler: func() *BaseHandler {
				mockPromoDB := repository.NewEmptyPromoDB()
				mockPromoDB.InsertPromoCode("VALIDCODE")
				productDB := repository.NewProductDB()

				return &BaseHandler{
					productDB: productDB,
					promoDB:   mockPromoDB,
					validator: validator.New(),
				}
			},
			expectedStatus: http.StatusCreated,
			expectedBody: `{
				"id": "mock-uuid",
				"items": [{"productId":"1","quantity":1},{"productId":"2","quantity":2}],
				"products": [{"id":"1","name":"Waffle with Berries","category":"Waffle","price":6.5},{"id":"2","name":"Vanilla Bean Crème Brûlée","category":"Crème Brûlée","price":7}]
			}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			decimal.MarshalJSONWithoutQuotes = true
			generateUUID = func() string {
				return "mock-uuid"
			}
			handler := tt.setupHandler()

			reqBody, err := json.Marshal(tt.orderRequest)
			assert.NoError(t, err)

			req, err := http.NewRequest("POST", "/order", bytes.NewReader(reqBody))
			assert.NoError(t, err)

			rr := httptest.NewRecorder()
			router := mux.NewRouter()
			router.HandleFunc("/order", handler.PlaceOrder)
			router.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			assert.JSONEq(t, tt.expectedBody, rr.Body.String())
		})
	}
}
