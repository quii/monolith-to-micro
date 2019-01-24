package cookme

// IngredientsRepo returns a collection of ingredients
type IngredientsRepo interface {
	Ingredients() PerishableIngredients
}

// IngredientsRepoFunc allows you to implement IngredientsRepo with a func
type IngredientsRepoFunc func() PerishableIngredients

// Ingredients returns the ingredients generated from f
func (f IngredientsRepoFunc) Ingredients() PerishableIngredients {
	return f()
}

// RecipeRepo returns a collection of recipes
type RecipeRepo interface {
	Recipes() Recipes
}

// RecipeRepoFunc allows you to implement RecipeRepo with a func
type RecipeRepoFunc func() Recipes

// Recipes returns recipes generated from f
func (f RecipeRepoFunc) Recipes() Recipes {
	return f()
}

// ListRecipes describes what meals should be cooked given the expiration dates of the IngredientsRepo
func ListRecipes(ingredientsRepo IngredientsRepo, recipeRepo RecipeRepo) Recipes {

	var foundRecipes Recipes
	ingredients := ingredientsRepo.Ingredients().SortByExpirationDate()

	for _, recipe := range recipeRepo.Recipes() {
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

	return foundRecipes
}
