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
