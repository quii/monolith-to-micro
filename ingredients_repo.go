package cookme

import (
	"github.com/google/go-cmp/cmp"
	"testing"
)

// IngredientsRepo returns a collection of ingredients
type IngredientsRepo interface {
	Ingredients() Ingredients
}

// IngredientsRepoFunc allows you to implement IngredientsRepo with a func
type IngredientsRepoFunc func() Ingredients

// Ingredients returns the ingredients generated from f
func (f IngredientsRepoFunc) Ingredients() Ingredients {
	return f()
}

// AssertIngredientsEqual is a test helper for checking if 2 lists of ingredients are the same
func AssertIngredientsEqual(t *testing.T, got, want Ingredients) {
	t.Helper()
	if !cmp.Equal(got, want) {
		t.Errorf("got %+v, want %+v", got, want)
	}
}
