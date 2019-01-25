package main

import (
	"github.com/quii/monolith-to-micro"
	"github.com/quii/monolith-to-micro/inventory"
	"github.com/quii/monolith-to-micro/recipe"
	"github.com/spf13/cobra"
	"log"
	"strconv"
	"time"
)

const dbFileName = "cookme.db"
const recipeAddress = "recipes:5000"

func main() {

	recipeBook, close := recipe.NewClient(recipeAddress)
	defer close()

	houseInventory, err := inventory.NewHouseInventory(dbFileName)

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

			log.Println("Why not cook")
			for _, recipe := range recipes {
				log.Println(recipe)
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

	var deleteIngredient = &cobra.Command{
		Use:   "delete-ingredient [name]",
		Short: "Delete ingredient from inventory",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			houseInventory.DeleteIngredient(args[0])
		},
	}

	var addRecipe = &cobra.Command{
		Use:   "add-recipe [name] [ingredients...]",
		Short: "Add recipe",
		Args:  cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			recipeBook.Add(args[0], args[1:])
		},
	}

	var deleteRecipe = &cobra.Command{
		Use:   "delete-recipe [name] ",
		Short: "Delete recipe",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			recipeBook.Delete(args[0])
		},
	}

	rootCmd.AddCommand(addIngredient)
	rootCmd.AddCommand(deleteIngredient)
	rootCmd.AddCommand(addRecipe)
	rootCmd.AddCommand(deleteRecipe)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
