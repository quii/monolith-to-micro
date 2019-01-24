package cookme

// DummyRecipeRepo returns a hard-coded list of recipes
func DummyRecipeRepo() Recipes {

	milk := Ingredient{Name: "Milk"}
	cheese := Ingredient{Name: "Cheese"}
	pasta := Ingredient{Name: "Pasta"}

	return Recipes{
		Recipe{Name: "Mac and cheese", Ingredients: Ingredients{pasta, cheese}},
		Recipe{Name: "Cheesy milk", Ingredients: Ingredients{milk, cheese}},
	}
}
