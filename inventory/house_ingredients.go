package inventory

import (
	"encoding/json"
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/quii/monolith-to-micro"
	"log"
	"time"
)

var (
	config   = &bolt.Options{Timeout: 1 * time.Second}
	bucket   = []byte("inventory")
	itemsKey = []byte("items")
)

type HouseInventory struct {
	dbFileName string
	inventory  cookme.Ingredients
}

func NewHouseInventory(dbFilename string) (*HouseInventory, error) {

	db, err := bolt.Open(dbFilename, 0600, config)

	if err != nil {
		return nil, fmt.Errorf("problem opening db '%s', %+v", dbFilename, err)
	}

	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(bucket)
		return err
	})

	return &HouseInventory{
		dbFileName: dbFilename,
	}, err
}

func (h *HouseInventory) Ingredients() cookme.Ingredients {
	db, err := bolt.Open(h.dbFileName, 0600, config)

	if err != nil {
		log.Printf("problem opening db '%s', %+v", h.dbFileName, err)
		return nil
	}

	defer db.Close()

	var ingredients cookme.Ingredients

	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)

		if err != nil {
			log.Println(err)
		}
		items := b.Get(itemsKey)

		if items != nil {
			return json.Unmarshal(items, &ingredients)
		}

		return nil
	})

	if err != nil {
		log.Printf("problem retrieving inventory %+v", err)
		return nil
	}

	return ingredients
}

func (h *HouseInventory) AddIngredients(ingredients ...cookme.Ingredient) {
	existingIngredients := h.Ingredients()

	newIngredients := append(existingIngredients, ingredients...)

	db, err := bolt.Open(h.dbFileName, 0600, config)

	if err != nil {
		log.Printf("problem opening db '%s', %+v", h.dbFileName, err)
		return
	}

	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)
		encodedIngredients, _ := json.Marshal(newIngredients)
		err := b.Put(itemsKey, encodedIngredients)

		if err != nil {
			log.Printf("problem adding ingredients to bucket %v", err)
			return nil
		}

		return nil
	})

}
