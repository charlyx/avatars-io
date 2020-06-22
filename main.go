package main

import (
	"log"
	"os"

	server "github.com/charlyx/avatars.io/server"
)

func main() {
	port := os.Getenv("PORT")

	log.Fatalf("Server stopped: %s", server.Start(port))
}
