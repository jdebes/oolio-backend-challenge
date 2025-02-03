package repository

import (
	"jdebes/oolio-backend-challenge/models"
)

type ProductDB struct {
	// Lookup map allows O(1) get access.
	lookup map[string]*models.Product
	store  []*models.Product
}

func NewProductDB() *ProductDB {
	return &ProductDB{
		lookup: make(map[string]*models.Product),
		store:  []*models.Product{},
	}
}

func (db *ProductDB) Insert(product models.Product) {
	if _, exists := db.lookup[product.ID]; !exists {
		db.store = append(db.store, &product)
	}
	db.lookup[product.ID] = &product
}

func (db *ProductDB) Get(productID string) (models.Product, bool) {
	product, exists := db.lookup[productID]
	if !exists {
		return models.Product{}, false
	}
	return *product, true
}

func (db *ProductDB) List() []models.Product {
	// Avoid accidentally changing internal state
	products := make([]models.Product, 0, len(db.store))
	for _, product := range db.store {
		products = append(products, *product)
	}
	return products
}
