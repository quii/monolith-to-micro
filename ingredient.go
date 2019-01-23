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

// Ingredient represents an ingredient and when it can be used by
type Ingredient struct {
	Name           string
	ExpirationDate time.Time
}

func (i Ingredient) String() string {
	expiresIn := math.Abs(math.Round(time.Since(i.ExpirationDate).Hours() / 24))
	return fmt.Sprintf("%s expires %v days", i.Name, expiresIn)
}

// Ingredients is a collection of Ingredient
type Ingredients []Ingredient

// Contains tells you if an ingredient exists in this slice
func (ingredients Ingredients) Contains(needle Ingredient) bool {
	for _, ingredient := range ingredients {
		if strings.ToLower(ingredient.Name) == strings.ToLower(needle.Name) {
			return true
		}
	}
	return false
}

// SortByExpirationDate sorts _in place_ the collection of ingredients
func (ingredients Ingredients) SortByExpirationDate() Ingredients {
	sort.Slice(ingredients, func(i, j int) bool {
		return ingredients[i].ExpirationDate.Before(ingredients[j].ExpirationDate)
	})

	return ingredients
}

// AssertIngredientsEqual is a test helper for checking if 2 lists of ingredients are the same
func AssertIngredientsEqual(t *testing.T, got, want Ingredients) {
	t.Helper()
	if !cmp.Equal(got, want) {
		t.Errorf("got %+v, want %+v", got, want)
	}
}
