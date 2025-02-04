package models

type Order struct {
	ID       string     `json:"id"`
	Items    []LineItem `json:"items"`
	Products []Product  `json:"products"`
}

type LineItem struct {
	ProductID string `json:"productId" validate:"required"`
	Quantity  int    `json:"quantity" validate:"required,gt=0"`
}

type OrderRequest struct {
	CouponCode string     `json:"couponCode" validate:"omitempty,min=8,max=10"`
	Items      []LineItem `json:"items" validate:"required,min=1,dive,required"`
}
