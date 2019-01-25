package inventory

import (
	"encoding/json"
	"github.com/quii/monolith-to-micro"
	"github.com/quii/monolith-to-micro/bucket"
	"log"
)

// HouseInventory manages PerishableIngredients, persisting the data in the filesystem
type HouseInventory struct {
	boltBucket *bucket.BoltBucket
}

const bucketName = "inventory"

// NewHouseInventory creates a new house inventory, creating the db file if needed
func NewHouseInventory(dbFilename string) (*HouseInventory, error) {
	bucket, err := bucket.NewBoltBucket(dbFilename, bucketName)

	inventory := &HouseInventory{
		boltBucket: bucket,
	}

	return inventory, err
}

// Ingredients lists all the ingredients in the house
func (h *HouseInventory) Ingredients() cookme.PerishableIngredients {
	var ingredients cookme.PerishableIngredients

	data, err := h.boltBucket.Get()

	if err != nil {
		log.Printf("problem getting data %+v", err)
		return nil
	}

	err = json.Unmarshal(data, &ingredients)

	return ingredients
}

// AddIngredients adds an ingredient to the inventory
func (h *HouseInventory) AddIngredients(ingredients ...cookme.PerishableIngredient) {
	existingIngredients := h.Ingredients()

	newIngredients := append(existingIngredients, ingredients...)

	h.boltBucket.Put(asJSON(newIngredients))
}

// DeleteIngredient will attempt to remove an ingredient from the inventory
func (h *HouseInventory) DeleteIngredient(ingredient string) {
	var newIngredients cookme.PerishableIngredients

	for _, i := range h.Ingredients() {
		if i.Name != ingredient {
			newIngredients = append(newIngredients, i)
		}
	}

	h.boltBucket.Put(asJSON(newIngredients))
}

func asJSON(ingredients cookme.PerishableIngredients) []byte {
	b, _ := json.Marshal(ingredients)
	return b
}
