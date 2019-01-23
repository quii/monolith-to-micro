package cookme_test

import (
	"bytes"
	"fmt"
	"github.com/quii/monolith-to-micro"
	"testing"
	"time"
)

func TestListIngredients(t *testing.T) {

	var someIngredients = []cookme.Ingredient{
		{Name: "Milk", ExpirationDate: time.Now().Add(72 * time.Hour)},
		{Name: "Cheese", ExpirationDate: time.Now().Add(48 * time.Hour)},
	}

	var StubIngredientsRepo = func() cookme.Ingredients {
		return someIngredients
	}

	t.Run("prints ingredients from ingredients list ordered by exp date", func(t *testing.T) {
		var got bytes.Buffer
		cookme.ListIngredients(&got, cookme.IngredientsRepoFunc(StubIngredientsRepo))

		want := fmt.Sprintf("%s\n%s\n", someIngredients[0], someIngredients[1])

		if got.String() != want {
			t.Errorf(`got "%s", want "%s"`, got.String(), want)
		}
	})
}
