package cookme

import "time"

func DummyIngredientsRepo() Ingredients {
	return []Ingredient{
		{"Cheese", time.Now().Add(48 * time.Hour)},
		{"Milk", time.Now().Add(72 * time.Hour)},
	}
}
