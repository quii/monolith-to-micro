package main

import (
	"fmt"
	"github.com/quii/monolith-to-micro"
	"github.com/quii/monolith-to-micro/inventory"
	"github.com/spf13/cobra"
	"log"
	"os"
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

			cookme.ListIngredients(
				os.Stdout,
				houseInventory,
			)
		},
	}

	var addIngredient = &cobra.Command{
		Use:   "add-ingredient",
		Short: "Add ingredient to inventory",
		Run: func(cmd *cobra.Command, args []string) {
			var newIngredients cookme.Ingredients
			for _, name := range args {
				newIngredients = append(newIngredients, cookme.Ingredient{name, time.Now().Add(72 * time.Hour)})
			}
			houseInventory.AddIngredients(newIngredients...)
		},
	}

	rootCmd.AddCommand(addIngredient)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
