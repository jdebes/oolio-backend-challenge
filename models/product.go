package models

import (
	"github.com/shopspring/decimal"
)

type Product struct {
	ID       string          `json:"id"`
	Name     string          `json:"name"`
	Price    decimal.Decimal `json:"price"`
	Category string          `json:"category"`
}
