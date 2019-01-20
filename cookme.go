package cookme

import (
	"fmt"
	"io"
)

func ListIngredients(out io.Writer, ingredientsRepo IngredientsRepo) {
	for _, ingredient := range ingredientsRepo.Ingredients().SortByExpirationDate() {
		fmt.Fprintln(out, ingredient)
	}
}
