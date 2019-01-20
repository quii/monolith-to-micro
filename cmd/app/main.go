package main

import (
	"github.com/quii/monolith-to-micro"
	"os"
)

func main() {
	cookme.ListIngredients(
		os.Stdout,
		cookme.IngredientsRepoFunc(cookme.DummyIngredientsRepo),
	)
}
