# Monolith to Micro

Starting a project with microservices I feel is a waste of time and a violation of YAGNI. That's not to say just build a big monolith! Begin with a "monolith" to start building your product, iterating and evolving it listening to feedback. Eventually you'll start to "see" the services you'll want to split out, based on real experience of the system rather than architect trying to guess it. 

From this project I hope we will see how you can successfully evolve toward a distributed system from a monolith so long as you maintain a decent level of refactoring. 

## To run

Assuming you have docker-compose installed

`docker-compose up`

## General ideas

- To keep running things consistent use docker-compose, even for the first iteration. That's not too much overhead and will make things gentler as we start to make our system distributed.
- Use gRPC to split things out.
- Keep it as a command line app just to minimise things.

## The problem

We want to know what to make for dinner!

There will be some kind of idea of what ingredients are in the house and what their expiration dates are. We'll also have a recipe book to derive meals from.

### How to break the problem down

1. **Hello, world**. Command line app running through docker-compose that prints hello, world
2. **Hard-coded ingredients to use**. Print out a list of ingredients that are available ordered by the expiration date
3. **Manage ingredients (add, delete)**
4. **Find meals** from a hardcoded list of recipes and print them instead, based on available ingredients
6. **Manage recipes (add, delete)**

At this point, we'll think about splitting into different gRPC services

### Possible further steps

1. **Return meals that dont have all ingredients** and list them
2. **Better ingredient management with quantities** so for example users can buy more eggs and add them in 

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

### Step 4

The default behaviour of the app right now is to `cookme.ListIngredients`. Let's rename that to `ListRecipes` to change our intent and then expand upon our existing tests to change the behaviour.

`ListRecipes` will need to depend on some kind of `RecipeRepo` to get recipes.

```go
type RecipeRepo interface {
	Recipes() Recipes
}

func ListRecipes(out io.Writer, ingredientsRepo IngredientsRepo, recipeRepo RecipeRepo) {
	//...
}
```

We'll need to update our existing tests to reflect the new wanted behaviour. 

```go
t.Run("prints recipes that can be cooked given the current ingredients", func(t *testing.T) {
    got := cookme.ListRecipes(
        newStubIngredientsRepo(milk, cheese, pasta),
        newStubRecipeRepo(macAndCheese, cheesyMilk),
    )

    want := cookme.Recipes{macAndCheese, cheesyMilk}

    cookme.AssertRecipesEqual(t, got, want)
})

t.Run("prints no recipes if there aren't any", func(t *testing.T) {
    got := cookme.ListRecipes(
        newStubIngredientsRepo(milk),
        newStubRecipeRepo(macAndCheese),
    )

    cookme.AssertRecipesEqual(t, got, nil)
})
```

Some notes:

- `ListRecipes` no longer takes an `io.Writer` to send its output as it was becoming unwieldy to test with and on retrospect doesn't seem like a good fit for the function. Instead it returns `Recipes` which are app can print out. 
- We have not prioritised based on when ingredients are going to expire yet.

As we have taken this MVP approach we can run the software and see that it _basically_ works with a hardcoded recipe book. Like last time the next step is to allow the user to manage recipes. 

When writing the tests it became clear our way of modelling ingredients isn't correct. We tried sharing the concept of `Ingredient` but it includes an expiration date which isn't relevant for recipes. 

We need to decouple an ingredient from the idea of it being perishable and then update our `inventory` and `recipe` packages to use the correct types. 

```go
type Ingredient struct {
	Name           string
}

type PerishableIngredient struct {
	Ingredient
	ExpirationDate time.Time
}
``` 

Here are the tests for our new persistent recipe book

```go
func TestRecipeBook(t *testing.T) {

	milk := cookme.Ingredient{Name: "Milk"}
	cheese := cookme.Ingredient{Name: "Cheese"}
	pasta := cookme.Ingredient{Name: "Pasta"}

	macAndCheese := cookme.NewRecipe("Mac and cheese", pasta, cheese)
	cheesyMilk := cookme.NewRecipe("Cheesy milk", milk, cheese)

	t.Run("returns no recipes when none have been added", func(t *testing.T) {
		book, cleanup := NewTestRecipeBook(t)
		defer cleanup()

		AssertRecipesEqual(t, book.Recipes(), nil)
	})

	t.Run("returns recipes when added", func(t *testing.T) {
		book, cleanup := NewTestRecipeBook(t)
		defer cleanup()

		book.Add(macAndCheese)
		book.Add(cheesyMilk)

		want := cookme.Recipes{macAndCheese, cheesyMilk}
		got := book.Recipes()

		AssertRecipesEqual(t, got, want)
	})

	t.Run("doesnt return recipes when deleted", func(t *testing.T) {
		book, cleanup := NewTestRecipeBook(t)
		defer cleanup()

		book.Add(macAndCheese)
		book.Add(cheesyMilk)
		book.Delete(macAndCheese)

		want := cookme.Recipes{cheesyMilk}
		got := book.Recipes()

		AssertRecipesEqual(t, got, want)
	})
}
```

We can integrate it in the same way in our application and if you now try it out we can now add recipes that will be listed to be cooked if you have the required recipes. 

## From monolith to microservices - The fun part!

We have our very basic MVP finished. Users can manage ingredients and recipes and our system will combine them to figure out what can be cooked. 

We've tried to keep the code reasonably abstracted. Let's take a look at our `ListRecipes`.

```go
func ListRecipes(ingredientsRepo IngredientsRepo, recipeRepo RecipeRepo) Recipes
```

It requires two dependencies via interfaces to do the job. This is great because it does not need to care about _how_ ingredients and recipes are delivered. Right now it is an in-process method call and both of them just fetch some recipes from disk. 

Importantly we didn't derive this from hours of discussions around a whiteboard, we arrived at it from _something real_. We iterated on our code and learned what abstractions we needed. If they were wrong we can change them very easily compared to changing a distributed system.

This is fine for now but lets pretend we want to make our recipe retrieving more sophisticated. Maybe it will have some ways of trawling recipe websites to scrape new recipes, maybe it's backed by an actual database and maybe it's doing a lot of complicated work and we will need to scale it horizontally. 

What we hope here is that gRPC can help us smoothly evolve our architecture.

### Define our recipe protocol

We will write a protobuf file which defines our recipe service based on the interface we've discovered through writing real software.

```proto
syntax = "proto3";

message Ingredient {
    string Name = 1;
}

message Recipe {
    string Name = 1;
    repeated Ingredient Ingredients = 2;
}

message GetRecipesRequest {
}

message GetRecipesResponse {
    repeated Recipe Recipes = 1;
}

service RecipeService {
    rpc GetRecipes (GetRecipesRequest) returns (GetRecipesResponse);
}
```

We'll just convert one of the methods to a gRPC call for now just to get the scaffolding together.

From there we can use the `protoc` command to generate Go code for clients and servers of this service. 

We will need a new application to run our recipe server so inside `cmd` we create a new folder `recipe` with a `main.go`. 

```go
package main

import (
	"github.com/quii/monolith-to-micro/recipe"
	"google.golang.org/grpc"
	"log"
	"net"
)

const dbFileName = "cookme.db"
const port = ":5000"

func main() {
	recipeBook, err := recipe.NewBook(dbFileName)
	
	listener, err := net.Listen("tcp", port)

	if err != nil {
		log.Fatalf("problem listening to port %s, %v", port, err)
	}

	server := grpc.NewServer()

	recipe.RegisterRecipeServiceServer(
		server,
		recipeBook,
	)

	if err := server.Serve(listener); err != nil {
		log.Fatalf("failed to serve %v", err)
	}
}
```

What's going on here?

- We create our `recipe.Book` as normal, using the file system as our store. 
- We need our new server to listen on a port which is done with `net.Listen`.
- We create a server with `grpc.NewServer` which takes care of all of the details of running a gRPC server.
- When we generated our code from the proto file we got a function called `RegisterRecipeServiceServer` which lets us... well register a recipe service with the gRPC server we created. 
- The second argument needs to implement the interface `RecipeServiceServer` which was also auto-generated for us.

Here is how our `recipe.Book` implements the interface. 

```go
func (b *Book) GetRecipes(c context.Context, r *GetRecipesRequest) (*GetRecipesResponse, error) {
	var recipes []*Recipe

	for _, r := range b.Recipes() {
		recipes = append(recipes, convertRecipeToGRPC(r))
	}

	return &GetRecipesResponse{Recipes: recipes}, nil
}
```

It's a bit wonky as the package `cookme` defines `Recipe` and now our protobuf version of it does too so our code has to convert between the two types. Other than that you can see it is very trivial to make our `Book` become a gRPC service.

Next we need to make it so our original application can connect to our new recipe server. 

gRPC has generated a client for our server `RecipeServiceClient` and we can very easily connect to our service

```go
conn, err := grpc.Dial(recipeAddress, grpc.WithInsecure())

if err != nil {
    log.Fatalf("could not connect to %s, %v", recipeAddress, err)
}

defer conn.Close()

recipeClient := recipe.NewRecipeServiceClient(conn)
```

And we can make our RPC call to get the recipes

```go
res, err := recipeClient.GetRecipes(context.Background(), &recipe.GetRecipesRequest{})
```

I encapsulated all the code into a type inside the `recipe` package so anyone with an address to a server can fetch recipes

```go
type Client struct {
	c RecipeServiceClient
}

func NewClient(address string) (client *Client, close func() error) {
	conn, err := grpc.Dial(address, grpc.WithInsecure())

	if err != nil {
		log.Fatalf("could not connect to %s, %v", address, err)
	}

	recipeClient := NewRecipeServiceClient(conn)

	return &Client{c: recipeClient}, conn.Close
}

func (c *Client) Recipes() cookme.Recipes {
	res, err := c.c.GetRecipes(context.Background(), &GetRecipesRequest{})

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
```

We need to update our `docker-compose` file so that both our applications are started and linked

```yaml
version: "3"

services:
  app:
    image: golang:1.11.5-alpine
    volumes:
      - .:/go/src/github.com/quii/monolith-to-micro
    working_dir: /go/src/github.com/quii/monolith-to-micro/cmd/app
    command: go run main.go
    links:
      - recipes

  recipes:
    image: golang:1.11.5-alpine
    volumes:
      - .:/go/src/github.com/quii/monolith-to-micro
    working_dir: /go/src/github.com/quii/monolith-to-micro/cmd/recipe
    command: go run main.go
    ports:
      - "5000"
```

Finally we just need to update our application to use our new client.

```go
const recipeAddress = "recipes:5000"

func main() {

	recipeBook, close := recipe.NewClient(recipeAddress)
	defer close()
	
	// later on use this recipeBook when calling cookme.ListRecipes
```


Hopefully you'll agree that to go from our "monolith" to a distributed version was relatively hassle free thanks to gRPC generating almost all of the required code, with us just having to write some boilerplate code to wire it together. 

gRPC gives us a number of benefits over a traditional "REST"ful approach

- No need to generate clients or servers; it's all derived from our proto files.
- Typesafe out of the box
- Protobuf messages are much smaller compared to JSON for network calls
- HTTP2 rather than inefficient HTTP1.1 (notice how that detail is entirely abstracted from us too, you wouldnt know if I hadn't told you!)
- Versioned out of the box. No more bikeshedding meetings about whether to do it in the URLs or in `Accept` headers!

We have some technical debt because not all of our method calls for recipes (add/delete) are available over the network yet. Let's see what it's like to change our `recipe.proto` and see how well the tooling helps us. 

Add the following to the proto file so we can add recipes over gRPC.

```proto
message AddRecipeRequest {
    Recipe Recipe = 1;
}

message AddRecipeResponse {
}

service RecipeService {
    rpc GetRecipes (GetRecipesRequest) returns (GetRecipesResponse);
    rpc AddRecipe (AddRecipeRequest) returns (AddRecipeResponse);
}
```

By running `./build.sh` it will regenerate the code and run our tests

`cmd/recipe/main.go:24:36: cannot use recipeBook (type *recipe.Book) as type recipe.RecipeServiceServer in argument to recipe.RegisterRecipeServiceServer:
 	*recipe.Book does not implement recipe.RecipeServiceServer (missing AddRecipe method)`
 	
The interface we need to implement to make our gRPC server has changed so we need to add the method. From there we can update our `Client` wrapper to expose the call and then update main. We can repeat this process for delete. 

What's great about this is the generated code _tells us what to write_ so we dont have to worry too much about the details and we get our system distributed. 
