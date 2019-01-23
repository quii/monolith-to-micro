package cookme

import (
	"fmt"
	"io"
)

// ListIngredients describes what ingredients should be used up ordered by expiration
func ListIngredients(out io.Writer, ingredientsRepo IngredientsRepo) {
	for _, ingredient := range ingredientsRepo.Ingredients().SortByExpirationDate() {
		fmt.Fprintln(out, ingredient)
	}
}
