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

			cookme.ListIngredients(
				os.Stdout,
				houseInventory,
			)
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

			newIngredient := cookme.Ingredient{
				Name:           args[0],
				ExpirationDate: time.Now().Add(time.Duration(daysExpire) * time.Hour),
			}
			houseInventory.AddIngredients(newIngredient)
		},
	}

	rootCmd.AddCommand(addIngredient)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
