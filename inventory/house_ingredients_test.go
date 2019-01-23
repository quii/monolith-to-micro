package inventory_test

import (
	"github.com/quii/monolith-to-micro"
	"github.com/quii/monolith-to-micro/inventory"
	"log"
	"os"
	"testing"
	"time"
)

func TestHouseInventory(t *testing.T) {

	milk := cookme.Ingredient{Name: "Milk", ExpirationDate: time.Now().Add(72 * time.Hour)}
	cheese := cookme.Ingredient{Name: "Cheese", ExpirationDate: time.Now().Add(48 * time.Hour)}

	t.Run("empty inventory returns no ingredients", func(t *testing.T) {
		inv, cleanup := NewTestInventory(t)
		defer cleanup()

		cookme.AssertIngredientsEqual(t, inv.Ingredients(), nil)
	})

	t.Run("adding an ingredient means it gets returned", func(t *testing.T) {
		inv, cleanup := NewTestInventory(t)
		defer cleanup()

		inv.AddIngredients(milk, cheese)

		cookme.AssertIngredientsEqual(t, inv.Ingredients(), cookme.Ingredients{milk, cheese})
	})

	t.Run("deleting an ingredient means it no longer gets returned", func(t *testing.T) {
		inv, cleanup := NewTestInventory(t)
		defer cleanup()

		inv.AddIngredients(milk, cheese)
		inv.DeleteIngredient(milk.Name)

		cookme.AssertIngredientsEqual(t, inv.Ingredients(), cookme.Ingredients{cheese})
	})
}

func NewTestInventory(t *testing.T) (inv *inventory.HouseInventory, cleanup func()) {
	t.Helper()
	dbFilename := RandomString() + ".db"
	inv, err := inventory.NewHouseInventory(dbFilename)

	if err != nil {
		log.Fatalf("problem creating inventory %+v", err)
	}

	return inv, func() {
		os.Remove(dbFilename)
	}
}
