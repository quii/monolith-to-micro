package cookme

import (
	"fmt"
	"math"
	"sort"
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

// SortByExpirationDate sorts _in place_ the collection of ingredients
func (ingredients Ingredients) SortByExpirationDate() Ingredients {
	sort.Slice(ingredients, func(i, j int) bool {
		return ingredients[i].ExpirationDate.Before(ingredients[j].ExpirationDate)
	})

	return ingredients
}
