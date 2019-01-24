package cookme

import (
	"github.com/google/go-cmp/cmp"
	"testing"
)

// Recipe represents a recipe with its required ingredients
type Recipe struct {
	Name        string
	Ingredients Ingredients
}

// NewRecipe is creates a recipe with some ingredients
func NewRecipe(name string, ingredients ...Ingredient) Recipe {
	return Recipe{Name: name, Ingredients: ingredients}
}

func (r Recipe) String() string {
	return r.Name
}

// Recipes is a slice of recipes
type Recipes []Recipe

// AssertRecipesEqual is a test helper for checking if 2 lists of recipes are the same
func AssertRecipesEqual(t *testing.T, got, want Recipes) {
	t.Helper()
	if !cmp.Equal(got, want) {
		t.Errorf("got %+v, want %+v", got, want)
	}
}
