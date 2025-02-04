package handlers

import (
	"github.com/go-playground/validator/v10"
	"jdebes/oolio-backend-challenge/repository"
)

type BaseHandler struct {
	productDB *repository.ProductDB
	promoDB   *repository.PromoDB
	validator *validator.Validate
}

func NewBaseHandler() *BaseHandler {
	productDB := repository.NewProductDB()
	promoDB := repository.NewPromoDB()

	return &BaseHandler{productDB: productDB, promoDB: promoDB, validator: validator.New()}
}
