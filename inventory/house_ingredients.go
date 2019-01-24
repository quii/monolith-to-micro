package inventory

import (
	"encoding/json"
	"github.com/quii/monolith-to-micro"
	"log"
)

// HouseInventory manages Ingredients, persisting the data in the filesystem
type HouseInventory struct {
	boltBucket *cookme.BoltBucket
}

const bucketName = "inventory"

// NewHouseInventory creates a new house inventory, creating the db file if needed
func NewHouseInventory(dbFilename string) (*HouseInventory, error) {
	bucket, err := cookme.NewBoltBucket(dbFilename, bucketName)

	inventory := &HouseInventory{
		boltBucket: bucket,
	}

	return inventory, err
}

// Ingredients lists all the ingredients in the house
func (h *HouseInventory) Ingredients() cookme.Ingredients {
	var ingredients cookme.Ingredients

	data, err := h.boltBucket.Get()

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

	h.boltBucket.Put(asJSON(newIngredients))
}

// DeleteIngredient will attempt to remove an ingredient from the inventory
func (h *HouseInventory) DeleteIngredient(ingredient string) {
	var newIngredients cookme.Ingredients

	for _, i := range h.Ingredients() {
		if i.Name != ingredient {
			newIngredients = append(newIngredients, i)
		}
	}

	h.boltBucket.Put(asJSON(newIngredients))
}

func asJSON(ingredients cookme.Ingredients) []byte {
	b, _ := json.Marshal(ingredients)
	return b
}
