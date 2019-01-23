package cookme_test

import (
	"github.com/quii/monolith-to-micro"
	"testing"
	"time"
)

func TestListIngredients(t *testing.T) {

	milk := cookme.Ingredient{Name: "Milk", ExpirationDate: time.Now().Add(72 * time.Hour)}
	cheese := cookme.Ingredient{Name: "Cheese", ExpirationDate: time.Now().Add(48 * time.Hour)}
	pasta := cookme.Ingredient{Name: "Pasta", ExpirationDate: time.Now().Add(2000 * time.Hour)}

	macAndCheese := cookme.Recipe{Name: "Mac and cheese", Ingredients: cookme.Ingredients{pasta, cheese}}
	cheesyMilk := cookme.Recipe{Name: "Cheesy milk", Ingredients: cookme.Ingredients{milk, cheese}}

	t.Run("prints recipes that can be cooked given the current ingredients", func(t *testing.T) {
		got := cookme.ListRecipes(
			newStubIngredientsRepo(milk, cheese, pasta),
			newStubRecipeRepo(macAndCheese, cheesyMilk),
		)

		want := cookme.Recipes{macAndCheese, cheesyMilk}

		cookme.AssertRecipesEqual(t, got, want)
	})

	t.Run("prints no recipes if there aren't any", func(t *testing.T) {
		got := cookme.ListRecipes(
			newStubIngredientsRepo(milk),
			newStubRecipeRepo(macAndCheese),
		)

		cookme.AssertRecipesEqual(t, got, nil)
	})
}

type stubIngredientsRepo struct {
	ingredients cookme.Ingredients
}

func newStubIngredientsRepo(ingredients ...cookme.Ingredient) *stubIngredientsRepo {
	return &stubIngredientsRepo{ingredients: ingredients}
}

func (s *stubIngredientsRepo) Ingredients() cookme.Ingredients {
	return s.ingredients
}

type stubRecipeRepo struct {
	recipes cookme.Recipes
}

func newStubRecipeRepo(recipes ...cookme.Recipe) *stubRecipeRepo {
	return &stubRecipeRepo{recipes: recipes}
}

func (s *stubRecipeRepo) Recipes() cookme.Recipes {
	return s.recipes
}
