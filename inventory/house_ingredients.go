package inventory

import "github.com/quii/monolith-to-micro"

type HouseInventory struct {
	inventory cookme.Ingredients
}

func NewHouseInventory() *HouseInventory {
	return &HouseInventory{}
}

func (h *HouseInventory) Ingredients() cookme.Ingredients {
	return h.inventory
}

func (h *HouseInventory) AddIngredients(ingredients ...cookme.Ingredient) {
	h.inventory = append(h.inventory, ingredients...)
}
