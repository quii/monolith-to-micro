package inventory

import (
	"encoding/json"
	"github.com/quii/monolith-to-micro"
	"log"
)

// HouseInventory manages Ingredients, persisting the data in the filesystem
type HouseInventory struct {
	boltBucket *boltBucket
}

// NewHouseInventory creates a new house inventory, creating the db file if needed
func NewHouseInventory(dbFilename string) (*HouseInventory, error) {
	bucket, err := newBoltBucket(dbFilename)

	inventory := &HouseInventory{
		boltBucket: bucket,
	}

	return inventory, err
}

// Ingredients lists all the ingredients in the house
func (h *HouseInventory) Ingredients() cookme.Ingredients {
	var ingredients cookme.Ingredients

	data, err := h.boltBucket.get()

	if err != nil {
		log.Printf("problem getting data %+v", err)
		return nil
	}

	err = json.Unmarshal(data, &ingredients)

	return ingredients
}

// AddIngredients adds an ingredient to the inventory
func (h *HouseInventory) AddIngredients(ingredients ...cookme.Ingredient) {
	existingIngredients := h.Ingredients()

	newIngredients := append(existingIngredients, ingredients...)

	encodedIngredients, _ := json.Marshal(newIngredients)

	h.boltBucket.put(encodedIngredients)
}
