package cookme

// FindRecipes finds appropriate recipes to cook given a list of recipes and perishable ingredients
func FindRecipes(recipes Recipes, ingredients PerishableIngredients) (foundRecipes Recipes) {
	for _, recipe := range recipes {
		allIngredientsFound := true
		for _, requiredIngredient := range recipe.Ingredients {
			if !ingredients.Contains(requiredIngredient) {
				allIngredientsFound = false
			}
		}

		if allIngredientsFound {
			foundRecipes = append(foundRecipes, recipe)
		}
	}

	return
}
