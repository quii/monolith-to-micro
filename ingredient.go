package cookme

import (
	"fmt"
	"github.com/google/go-cmp/cmp"
	"math"
	"sort"
	"strings"
	"testing"
	"time"
)

// Ingredient represents an ingredient for cooking
type Ingredient struct {
	Name string
}

// ExpiresAt returns a PerishableIngredient which expires at t
func (i Ingredient) ExpiresAt(t time.Time) PerishableIngredient {
	return PerishableIngredient{i, t}
}

// Ingredients is a collection of Ingredients
type Ingredients []Ingredient

// PerishableIngredient represents an ingredient and when it can be used by
type PerishableIngredient struct {
	Ingredient
	ExpirationDate time.Time
}

func (p PerishableIngredient) String() string {
	expiresIn := math.Abs(math.Round(time.Since(p.ExpirationDate).Hours() / 24))
	return fmt.Sprintf("%s expires %v days", p.Name, expiresIn)
}

// PerishableIngredients is a collection of PerishableIngredient
type PerishableIngredients []PerishableIngredient

// Contains tells you if an ingredient exists in this slice
func (ingredients PerishableIngredients) Contains(needle Ingredient) bool {
	for _, ingredient := range ingredients {
		if strings.ToLower(ingredient.Name) == strings.ToLower(needle.Name) {
			return true
		}
	}
	return false
}

// SortByExpirationDate sorts _in place_ the collection of ingredients
func (ingredients PerishableIngredients) SortByExpirationDate() PerishableIngredients {
	sort.Slice(ingredients, func(i, j int) bool {
		return ingredients[i].ExpirationDate.Before(ingredients[j].ExpirationDate)
	})

	return ingredients
}

// AssertPerishableIngredientsEqual is a test helper for checking if 2 lists of ingredients are the same
func AssertPerishableIngredientsEqual(t *testing.T, got, want PerishableIngredients) {
	t.Helper()
	if !cmp.Equal(got, want) {
		t.Errorf("got %+v, want %+v", got, want)
	}
}
