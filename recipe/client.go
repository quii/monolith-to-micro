package recipe

import (
	"context"
	"github.com/quii/monolith-to-micro"
	"google.golang.org/grpc"
	"log"
)

// Client is a RecipeRepo connecting to the recipe server
type Client struct {
	RecipeServiceClient
}

// NewClient creates a new client to the recipe server, make sure to call defer close()
func NewClient(address string) (client *Client, close func() error) {
	conn, err := grpc.Dial(address, grpc.WithInsecure())

	if err != nil {
		log.Fatalf("could not connect to %s, %v", address, err)
	}

	recipeClient := NewRecipeServiceClient(conn)

	return &Client{RecipeServiceClient: recipeClient}, conn.Close
}

// Recipes returns all recipes available
func (c *Client) Recipes() cookme.Recipes {
	request := &GetRecipesRequest{}
	res, err := c.GetRecipes(context.Background(), request)

	if err != nil {
		log.Fatalf("problem getting recipes %v", err)
	}

	var recipes cookme.Recipes

	for _, r := range res.Recipes {
		var ingredients cookme.Ingredients
		for _, i := range r.Ingredients {
			ingredients = append(ingredients, cookme.Ingredient{Name: i.Name})
		}
		recipes = append(recipes, cookme.Recipe{
			Name:        r.Name,
			Ingredients: ingredients,
		})
	}

	return recipes
}
