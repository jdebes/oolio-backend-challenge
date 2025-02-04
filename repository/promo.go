package repository

import (
	"bufio"
	"os"

	log "github.com/sirupsen/logrus"
	"jdebes/oolio-backend-challenge/data"
)

type PromoDB struct {
	db map[string]interface{}
}

func NewPromoDB() *PromoDB {
	promoStore, err := readPromos()
	if err != nil {
		log.Errorf("Error reading promos, will contine with empty store: %v.", err)
	}

	return &PromoDB{
		db: promoStore,
	}
}

func NewEmptyPromoDB() *PromoDB {
	return &PromoDB{
		db: make(map[string]interface{}),
	}
}

func (db *PromoDB) IsValidPromoCode(code string) bool {
	_, exists := db.db[code]
	return exists
}

func (db *PromoDB) InsertPromoCode(code string) {
	db.db[code] = struct{}{}
}

func readPromos() (map[string]interface{}, error) {
	promoStore := make(map[string]interface{})
	file, err := os.Open(data.DirPrefix + data.ValidPromosFile)
	if err != nil {
		return promoStore, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		promoCode := scanner.Text()
		promoStore[promoCode] = struct{}{}
	}

	return promoStore, nil
}
