package main

import (
	"fmt"
	"github.com/quii/monolith-to-micro"
	"github.com/quii/monolith-to-micro/inventory"
	"github.com/spf13/cobra"
	"log"
	"os"
	"strconv"
	"time"
)

func main() {
	houseInventory, err := inventory.NewHouseInventory("inventory.db")

	if err != nil {
		log.Fatalf("problem creating inventory %v", err)
	}

	var rootCmd = &cobra.Command{
		Use:   "cookme",
		Short: "Cook me tells you what you should cook",
		Run: func(cmd *cobra.Command, args []string) {

			recipes := cookme.ListRecipes(
				houseInventory,
				cookme.RecipeRepoFunc(cookme.DummyRecipeRepo),
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

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
