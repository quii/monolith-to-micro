package cookme

import (
	"github.com/google/go-cmp/cmp"
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
	if !cmp.Equal(got, want) {
		t.Errorf("got %+v, want %+v", got, want)
	}
}
