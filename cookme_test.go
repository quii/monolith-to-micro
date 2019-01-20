package cookme

import (
	"bytes"
	"fmt"
	"testing"
	"time"
)

func TestListIngredients(t *testing.T) {

	var someIngredients = []Ingredient{
		{"Milk", time.Now().Add(72 * time.Hour)},
		{"Cheese", time.Now().Add(48 * time.Hour)},
	}

	var StubIngredientsRepo = func() Ingredients {
		return someIngredients
	}

	t.Run("prints ingredients from ingredients list ordered by exp date", func(t *testing.T) {
		var got bytes.Buffer
		ListIngredients(&got, IngredientsRepoFunc(StubIngredientsRepo))

		want := fmt.Sprintf("%s\n%s\n", someIngredients[0], someIngredients[1])

		if got.String() != want {
			t.Errorf(`got "%s", want "%s"`, got.String(), want)
		}
	})
}
