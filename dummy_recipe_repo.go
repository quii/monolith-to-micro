package cookme

import "time"

// DummyRecipeRepo returns a hard-coded list of recipes
func DummyRecipeRepo() Recipes {

	milk := Ingredient{Name: "Milk", ExpirationDate: time.Now().Add(72 * time.Hour)}
	cheese := Ingredient{Name: "Cheese", ExpirationDate: time.Now().Add(48 * time.Hour)}
	pasta := Ingredient{Name: "Pasta", ExpirationDate: time.Now().Add(2000 * time.Hour)}

	return Recipes{
		Recipe{Name: "Mac and cheese", Ingredients: Ingredients{pasta, cheese}},
		Recipe{Name: "Cheesy milk", Ingredients: Ingredients{milk, cheese}},
	}
}
