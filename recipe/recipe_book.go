package recipe

import (
	"context"
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

// GetRecipes allows Book to act as a RecipeServiceServer
func (b *Book) GetRecipes(c context.Context, r *GetRecipesRequest) (*GetRecipesResponse, error) {
	var recipes []*Recipe

	for _, r := range b.Recipes() {
		recipes = append(recipes, convertRecipeToGRPC(r))
	}

	return &GetRecipesResponse{Recipes: recipes}, nil
}

// AddRecipe allows Book to act as a RecipeServiceServer
func (b *Book) AddRecipe(ctx context.Context, in *AddRecipeRequest) (*AddRecipeResponse, error) {
	newRecipe := cookme.Recipe{Name: in.Recipe.Name}

	for _, i := range in.Recipe.Ingredients {
		newRecipe.Ingredients = append(newRecipe.Ingredients, cookme.Ingredient{Name: i.Name})
	}

	b.Add(newRecipe)

	response := &AddRecipeResponse{}
	return response, nil
}

// Recipes returns all recipes
func (b *Book) Recipes() cookme.Recipes {
	var recipes cookme.Recipes
	stuff, _ := b.boltBucket.Get()
	json.Unmarshal(stuff, &recipes)
	return recipes
}

// Add will add a recipe to the book
func (b *Book) Add(recipe cookme.Recipe) {
	newRecipes := append(b.Recipes(), recipe)
	b.boltBucket.Put(asJSON(newRecipes))
}

// Delete will remove a recipe from the book
func (b *Book) Delete(name string) {
	var newRecipes cookme.Recipes

	for _, r := range b.Recipes() {
		if r.Name != name {
			newRecipes = append(newRecipes, r)
		}
	}

	b.boltBucket.Put(asJSON(newRecipes))
}

func asJSON(recipes cookme.Recipes) []byte {
	b, _ := json.Marshal(recipes)
	return b
}

func convertRecipeToGRPC(r cookme.Recipe) *Recipe {
	var ingredients []*Ingredient
	for _, i := range r.Ingredients {
		ingredients = append(ingredients, &Ingredient{Name: i.Name})
	}
	recipe := &Recipe{Name: r.Name, Ingredients: ingredients}
	return recipe
}
