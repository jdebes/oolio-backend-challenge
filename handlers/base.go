package handlers

import (
	"github.com/go-playground/validator/v10"
	"github.com/shopspring/decimal"
	"jdebes/oolio-backend-challenge/models"
	"jdebes/oolio-backend-challenge/repository"
)

type BaseHandler struct {
	db        *repository.ProductDB
	validator *validator.Validate
}

func NewBaseHandler() *BaseHandler {
	db := repository.NewProductDB()

	products := []models.Product{
		{ID: "1", Name: "Waffle with Berries", Category: "Waffle", Price: decimal.NewFromFloat(6.5)},
		{ID: "2", Name: "Vanilla Bean Crème Brûlée", Category: "Crème Brûlée", Price: decimal.NewFromFloat(7)},
		{ID: "3", Name: "Macaron Mix of Five", Category: "Macaron", Price: decimal.NewFromFloat(8)},
		{ID: "4", Name: "Classic Tiramisu", Category: "Tiramisu", Price: decimal.NewFromFloat(5.5)},
		{ID: "5", Name: "Pistachio Baklava", Category: "Baklava", Price: decimal.NewFromFloat(4)},
	}

	for _, product := range products {
		db.Insert(product)
	}

	return &BaseHandler{db: db, validator: validator.New()}
}
