package main

import (
	"fmt"
	"github.com/quii/monolith-to-micro"
	"github.com/quii/monolith-to-micro/inventory"
	"github.com/spf13/cobra"
	"os"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "cookme",
		Short: "Cook me tells you what you should cook",
		Run: func(cmd *cobra.Command, args []string) {
			cookme.ListIngredients(
				os.Stdout,
				inventory.NewHouseInventory(),
			)
		},
	}

	var addIngredient = &cobra.Command{
		Use:   "add-ingredient",
		Short: "Add ingredient to inventory",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("will add ingredients %+v\n", args)
		},
	}

	rootCmd.AddCommand(addIngredient)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
