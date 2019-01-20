package cookme

import (
	"fmt"
	"math"
	"sort"
	"time"
)

type Ingredient struct {
	Name string
	ExpirationDate time.Time
}

func (i Ingredient) String() string {
	expiresIn := math.Abs(math.Round(time.Since(i.ExpirationDate).Hours() / 24))
	return fmt.Sprintf("%s expires %v days", i.Name, expiresIn)
}

type Ingredients []Ingredient

func (ingredients Ingredients) SortByExpirationDate() Ingredients  {
	sort.Slice(ingredients, func(i, j int) bool {
		return ingredients[i].ExpirationDate.Before(ingredients[j].ExpirationDate)
	})

	return ingredients
}
