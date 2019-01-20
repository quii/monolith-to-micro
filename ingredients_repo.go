package cookme

type IngredientsRepo interface {
	Ingredients() Ingredients
}

type IngredientsRepoFunc func() Ingredients

func (f IngredientsRepoFunc) Ingredients() Ingredients {
	return f()
}
