package cookme

import "time"

var (
	twoDays   = time.Now().Add(48 * time.Hour)
	threeDays = time.Now().Add(72 * time.Hour)
)

func DummyIngredientsRepo() Ingredients {
	return []Ingredient{
		{"Cheese", twoDays},
		{"Milk", threeDays},
	}
}
