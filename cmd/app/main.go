package main

import (
	"fmt"
	"github.com/quii/monolith-to-micro"
	"github.com/quii/monolith-to-micro/inventory"
	"github.com/quii/monolith-to-micro/recipe"
	"github.com/spf13/cobra"
	"log"
	"os"
	"strconv"
	"time"
)

const dbFileName = "cookme.db"

func main() {
	houseInventory, err := inventory.NewHouseInventory(dbFileName)
	recipeBook, err := recipe.NewBook(dbFileName)

	if err != nil {
		log.Fatalf("problem creating db %v", err)
	}

	var rootCmd = &cobra.Command{
		Use:   "cookme",
		Short: "Cook me tells you what you should cook",
		Run: func(cmd *cobra.Command, args []string) {

			recipes := cookme.ListRecipes(
				houseInventory,
				recipeBook,
			)

			for _, recipe := range recipes {
				fmt.Println(recipe)
			}
		},
	}

	var addIngredient = &cobra.Command{
		Use:   "add-ingredient [name] [days-to-expire]",
		Short: "Add ingredient to inventory",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			hoursExpire, err := strconv.Atoi(args[1])

			if err != nil {
				log.Fatalf("invalid days argument, expect a number")
			}

			daysExpire := hoursExpire * 24

			newIngredient := cookme.PerishableIngredient{
				Ingredient:     cookme.Ingredient{Name: args[0]},
				ExpirationDate: time.Now().Add(time.Duration(daysExpire) * time.Hour),
			}
			houseInventory.AddIngredients(newIngredient)
		},
	}

	rootCmd.AddCommand(addIngredient)

	var deleteIngredient = &cobra.Command{
		Use:   "delete-ingredient [name]",
		Short: "Delete ingredient from inventory",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			houseInventory.DeleteIngredient(args[0])
		},
	}

	rootCmd.AddCommand(deleteIngredient)

	var addRecipe = &cobra.Command{
		Use:   "add-recipe [name] [ingredients...]",
		Short: "Add recipe",
		Args: cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			recipeName := args[0]
			var ingredients cookme.Ingredients
			for _, i := range args[1:] {
				ingredients = append(ingredients, cookme.Ingredient{Name: i})
			}
			recipeBook.Add(cookme.Recipe{Name: recipeName, Ingredients: ingredients})
		},
	}

	rootCmd.AddCommand(addRecipe)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
