package cookme

import "log"

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

	ingredients := ingredientsRepo.Ingredients().SortByExpirationDate()
	recipes := recipeRepo.Recipes()

	log.Printf("All ingredients %+v\n", ingredients)
	log.Printf("All recipes %+v\n", recipes)

	return FindRecipes(recipes, ingredients)
}
