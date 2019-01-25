package recipe

import (
	"encoding/json"
	"github.com/quii/monolith-to-micro"
	"github.com/quii/monolith-to-micro/bucket"
)

// Book contains recipes
type Book struct {
	boltBucket *bucket.BoltBucket
	recipes    cookme.Recipes
}

const boltBucketName = "recipes"

// NewBook returns a new recipe book backed by a bolt db at dbFilename
func NewBook(dbFilename string) (*Book, error) {
	boltBucket, err := bucket.NewBoltBucket(dbFilename, boltBucketName)

	if err != nil {
		return nil, err
	}

	return &Book{boltBucket: boltBucket}, nil
}

// Recipes returns all recipes
func (r Book) Recipes() cookme.Recipes {
	var recipes cookme.Recipes
	stuff, _ := r.boltBucket.Get()
	json.Unmarshal(stuff, &recipes)
	return recipes
}

// Add will add a recipe to the book
func (r *Book) Add(recipe cookme.Recipe) {
	newRecipes := append(r.Recipes(), recipe)
	r.boltBucket.Put(asJSON(newRecipes))
}

// Delete will remove a recipe from the book
func (r *Book) Delete(recipe cookme.Recipe) {
	var newRecipes cookme.Recipes

	for _, r := range r.Recipes() {
		if r.Name != recipe.Name {
			newRecipes = append(newRecipes, r)
		}
	}

	r.boltBucket.Put(asJSON(newRecipes))
}

func asJSON(recipes cookme.Recipes) []byte {
	b, _ := json.Marshal(recipes)
	return b
}
