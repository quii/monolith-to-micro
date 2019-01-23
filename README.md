# Monolith to Micro

Starting microservices I feel is a waste of time and a violation of YAGNI. That's not to say just build a big monolith! You should start with a monolith to start building your product, iterating and evolving it listening to feedback. Eventually you'll start to "see" the services you'll want to split out, based on real experience of the system rather than architect trying to guess it. 

From this project I hope we will see how you can successfully evolve toward a distributed system from a monolith so long as you maintain a decent level of refactoring. 

## To run

Assuming you have docker-compose installed

`docker-compose up`

## General ideas

- To keep running things consistent use docker-compose, even for the first iteration. That's not too much overhead and will make things gentler as we add new things.
- Use gRPC to split things out.
- Keep it as a command line app just to minimise things.

## The problem

We want to know what to make for dinner!

There will be some kind of idea of what ingredients are in the house and what their expiration dates are. We'll also have a recipe book to derive meals from, which eventually we should be able to filter by 

### How to break the problem down

1. **Hello, world**. Command line app running through docker-compose that prints hello, world
2. **Hard-coded ingredients to use**. Print out a list of ingredients that are available ordered by the expiration date
3. **Manage ingredients (add, delete)**
4. **Find meals** from a hardcoded list of recipes and print them instead, based on available ingredients
5. **Return meals that dont have all ingredients** and list them
6. **Manage ingredients**

At this point, we'll think about splitting into different gRPC services 

## Diary

### Step 1

As documented the goal of this iteration is to setup a simple hello world project. 

Created a `cookme.go` in the root of the project with one function 

```go
func ListIngredients(out io.Writer) {
	fmt.Fprintln(out, "Hello, world")
}
```

Then created a `/cmd/app` folder with a `main.go` which calls that function with `os.Stdout`. It's not over the top to separate our "library" code away from the app and this little bit of structure lets us setup docker-compose to run our app. 

```yaml
version: "3"

services:
  app:
    image: golang:1.11.2-alpine
    volumes:
      - .:/go/src/github.com/quii/monolith-to-micro
    working_dir: /go/src/github.com/quii/monolith-to-micro/cmd/app
    command: go run main.go
```

This `docker-compose.yaml` lives in our root and allows us to run our application in a container. This gives a common way of running our code and will become important later if we wish to add other dependencies such as databases or our own services if we evolve our architecture. 

### Step 2

Next we want to print out a list of ingredients from a hard-coded list. At this point it felt prudent to write a test for our `ListIngredients` function and we have extended it so it has a dependency of an `IngredientsRepo`

```go
func ListIngredients(out io.Writer, ingredientsRepo IngredientsRepo) {
	for _, ingredient := range ingredientsRepo.Ingredients().SortByExpirationDate() {
		fmt.Fprintln(out, ingredient)
	}
}
```

This test has revealed a potential abstraction to further build on in terms of _something_ to get ingredients from. At this stage though I have resisted the temptation of creating a new package (or service) as the code still feels manageable and at present the abstraction doesnt give us much right now.

### Step 3 

Next we're giving the ability for the user to manage the ingredients so we need a way to send commands to our application. For this we're using the excellent [Cobra](https://github.com/spf13/cobra) library.

```go
func main() {
	var rootCmd = &cobra.Command{
		Use:   "cookme",
		Short: "Cook me tells you what you should cook",
		Run: func(cmd *cobra.Command, args []string) {
			cookme.ListIngredients(
				os.Stdout,
				cookme.IngredientsRepoFunc(cookme.DummyIngredientsRepo),
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
```

We've added our new command `add-ingredient` and for now we just print a debug message. 

To run `add-ingredient` via docker-compose you will need to do `docker-compose run app go run main.go add-ingredient` (I am 100% sure this can be improved!)

We'll next need to add to our code a means of adding ingredients. This means we will have to move away from `DummyIngredientsRepo` which is just a hardcoded list of ingredients into something that can maintain state. 

There will actually be a fair amount of domain logic within this code and lots of tests. This feels like our `package cookme` may start to have too many concerns mixed with it so we'll start a new package called `inventory` and put a skeleton implementation of something that implements `IngredientsRepo` and gives us a function to add ingredients. 

```go
package inventory

import "github.com/quii/monolith-to-micro"

type HouseInventory struct {
	
}

func NewHouseInventory() *HouseInventory {
	return &HouseInventory{}
}

func (h *HouseInventory) Ingredients() cookme.Ingredients {
	panic("not implemented")
}

func (h *HouseInventory) AddIngredients(ingredients ...cookme.Ingredient) {
	panic(" not implemented")
}
```

Then replace the `DummyIngredientsRepo` with our new implementation in our application

```go
cookme.ListIngredients(
    os.Stdout,
    inventory.NewHouseInventory(),
)
```

If you try and run `docker-compose up` it should _compile_ but panic because we have not implemented our new code yet. We can drive this out with some tests.

```go
func TestHouseInventory(t *testing.T) {

	t.Run("empty inventory returns no ingredients", func(t *testing.T) {
		inv, cleanup := NewTestInventory(t)
		defer cleanup()

		cookme.AssertIngredientsEqual(t, inv.Ingredients(), nil)
	})

	t.Run("adding an ingredient means it gets returned", func(t *testing.T) {
		inv, cleanup := NewTestInventory(t)
		defer cleanup()

		milk := cookme.Ingredient{Name: "Milk", ExpirationDate: time.Now().Add(72 * time.Hour)}
		cheese := cookme.Ingredient{Name: "Cheese", ExpirationDate: time.Now().Add(48 * time.Hour)}

		inv.AddIngredients(milk, cheese)

		cookme.AssertIngredientsEqual(t, inv.Ingredients(), cookme.Ingredients{milk, cheese})
	})
}
```

You can check out the code that makes this pass in the repository but it's not especially interesting other than I decided to use [BoltDB](https://github.com/boltdb/bolt) to persist the inventory to disk, for fun. 

Now we have a working inventory we can update our application code to support adding ingredients properly

```go
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
```

We can carry on this approach for deleting ingredients. We'll add a test for our `HouseInventory` to add a new `Delete` method and then wire it up into our app. 

```go
t.Run("deleting an ingredient means it no longer gets returned", func(t *testing.T) {
    inv, cleanup := NewTestInventory(t)
    defer cleanup()

    inv.AddIngredients(milk, cheese)
    inv.DeleteIngredient(milk.Name)

    cookme.AssertIngredientsEqual(t, inv.Ingredients(), cookme.Ingredients{cheese})
})
```
```go
	var deleteIngredient = &cobra.Command{
		Use:   "delete-ingredient [name]",
		Short: "Delete ingredient from inventory",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			houseInventory.DeleteIngredient(args[0])
		},
	}

	rootCmd.AddCommand(deleteIngredient)
```

There's a number of short-comings with our software.

A lot of the time you dont "delete" an ingredient, you use _some_ of it. So at some point we will need to have the concept of ingredients having a quantity and when you add/delete the inventory will keep track of the totals.

However we have _working software_, it's MVP and it's not perfect but it's better for us to explore the broad ideas first so we get answers to the big questions about how to structure our app.  

The next important functionality to tackle is to have the software suggest what to cook given some recipes and the current state of the inventory. 