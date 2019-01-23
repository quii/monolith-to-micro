package inventory

import (
	"encoding/json"
	"github.com/quii/monolith-to-micro"
	"log"
)

type HouseInventory struct {
	boltBucket *boltBucket
}

func NewHouseInventory(dbFilename string) (*HouseInventory, error) {
	bucket, err := newBoltBucket(dbFilename)

	inventory := &HouseInventory{
		boltBucket: bucket,
	}

	return inventory, err
}

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

func (h *HouseInventory) AddIngredients(ingredients ...cookme.Ingredient) {
	existingIngredients := h.Ingredients()

	newIngredients := append(existingIngredients, ingredients...)

	encodedIngredients, _ := json.Marshal(newIngredients)

	h.boltBucket.put(encodedIngredients)
}
