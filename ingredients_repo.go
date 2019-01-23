package cookme

import (
	"reflect"
	"testing"
)

type IngredientsRepo interface {
	Ingredients() Ingredients
}

type IngredientsRepoFunc func() Ingredients

func (f IngredientsRepoFunc) Ingredients() Ingredients {
	return f()
}

func AssertIngredientsEqual(t *testing.T, got, want Ingredients) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %+v, want %+v", got, want)
	}
}
