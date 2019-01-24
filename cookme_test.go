package cookme_test

import (
	"github.com/quii/monolith-to-micro"
	"testing"
	"time"
)

func TestListIngredients(t *testing.T) {

	milk := cookme.Ingredient{Name: "Milk"}
	cheese := cookme.Ingredient{Name: "Cheese"}
	pasta := cookme.Ingredient{Name: "Pasta"}

	macAndCheese := cookme.Recipe{Name: "Mac and cheese", Ingredients: cookme.Ingredients{pasta, cheese}}
	cheesyMilk := cookme.Recipe{Name: "Cheesy milk", Ingredients: cookme.Ingredients{milk, cheese}}

	t.Run("prints recipes that can be cooked given the current ingredients", func(t *testing.T) {
		got := cookme.ListRecipes(
			newStubIngredientsRepo(
				milk.ExpiresAt(time.Now().Add(72*time.Hour)),
				cheese.ExpiresAt(time.Now().Add(48*time.Hour)),
				pasta.ExpiresAt(time.Now().Add(2000*time.Hour)),
			),
			newStubRecipeRepo(macAndCheese, cheesyMilk),
		)

		want := cookme.Recipes{macAndCheese, cheesyMilk}

		cookme.AssertRecipesEqual(t, got, want)
	})

	t.Run("prints no recipes if there aren't any", func(t *testing.T) {
		got := cookme.ListRecipes(
			newStubIngredientsRepo(milk.ExpiresAt(time.Now().Add(72*time.Hour))),
			newStubRecipeRepo(macAndCheese),
		)

		cookme.AssertRecipesEqual(t, got, nil)
	})
}

type stubIngredientsRepo struct {
	ingredients cookme.PerishableIngredients
}

func newStubIngredientsRepo(ingredients ...cookme.PerishableIngredient) *stubIngredientsRepo {
	return &stubIngredientsRepo{ingredients: ingredients}
}

func (s *stubIngredientsRepo) Ingredients() cookme.PerishableIngredients {
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
