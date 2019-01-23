package inventory_test

import (
	"github.com/quii/monolith-to-micro"
	"github.com/quii/monolith-to-micro/inventory"
	"testing"
	"time"
)

func TestHouseInventory(t *testing.T) {

	t.Run("empty inventory returns no ingredients", func(t *testing.T) {
		inv := inventory.NewHouseInventory()

		cookme.AssertIngredientsEqual(t, inv.Ingredients(), nil)
	})

	t.Run("adding an ingredient means it gets returned", func(t *testing.T) {
		inv := inventory.NewHouseInventory()
		milk := cookme.Ingredient{Name: "Milk", ExpirationDate: time.Now().Add(72 * time.Hour)}
		cheese := cookme.Ingredient{Name: "Cheese", ExpirationDate: time.Now().Add(48 * time.Hour)}

		inv.AddIngredients(milk, cheese)

		cookme.AssertIngredientsEqual(t, inv.Ingredients(), cookme.Ingredients{milk, cheese})
	})
}
