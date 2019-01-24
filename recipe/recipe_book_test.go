package recipe_test

import (
	"github.com/google/go-cmp/cmp"
	"github.com/quii/monolith-to-micro"
	"github.com/quii/monolith-to-micro/recipe"
	"log"
	"os"
	"testing"
	"time"
)

func TestRecipeBook(t *testing.T) {

	milk := cookme.Ingredient{Name: "Milk", ExpirationDate: time.Now().Add(72 * time.Hour)}
	cheese := cookme.Ingredient{Name: "Cheese", ExpirationDate: time.Now().Add(48 * time.Hour)}
	pasta := cookme.Ingredient{Name: "Pasta", ExpirationDate: time.Now().Add(2000 * time.Hour)}

	macAndCheese := cookme.Recipe{Name: "Mac and cheese", Ingredients: cookme.Ingredients{pasta, cheese}}
	cheesyMilk := cookme.Recipe{Name: "Cheesy milk", Ingredients: cookme.Ingredients{milk, cheese}}

	t.Run("returns no recipes when none have been added", func(t *testing.T) {
		book, cleanup := NewTestRecipeBook(t)
		defer cleanup()

		AssertRecipesEqual(t, book.Recipes(), nil)
	})

	t.Run("returns recipes when added", func(t *testing.T) {
		book, cleanup := NewTestRecipeBook(t)
		defer cleanup()

		book.Add(macAndCheese)
		book.Add(cheesyMilk)

		want := cookme.Recipes{macAndCheese, cheesyMilk}
		got := book.Recipes()

		AssertRecipesEqual(t, got, want)
	})

	t.Run("doesnt return recipes when deleted", func(t *testing.T) {
		book, cleanup := NewTestRecipeBook(t)
		defer cleanup()

		book.Add(macAndCheese)
		book.Add(cheesyMilk)
		book.Delete(macAndCheese)

		want := cookme.Recipes{cheesyMilk}
		got := book.Recipes()

		AssertRecipesEqual(t, got, want)
	})
}

func AssertRecipesEqual(t *testing.T, got, want cookme.Recipes) {
	t.Helper()
	if !cmp.Equal(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func NewTestRecipeBook(t *testing.T) (inv *recipe.Book, cleanup func()) {
	t.Helper()
	dbFilename := cookme.RandomString() + ".db"
	inv, err := recipe.NewBook(dbFilename)

	if err != nil {
		log.Fatalf("problem creating db %+v", err)
	}

	return inv, func() {
		os.Remove(dbFilename)
	}
}
