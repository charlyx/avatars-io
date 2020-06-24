package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/charlyx/avatars.io/secrets"
	"github.com/charlyx/avatars.io/server"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	secretClient, err := secrets.NewClient(os.Getenv("PROJECT_ID"))
	if err != nil {
		log.Fatalf("could not create secret accessor: %s", err)
	}

	handler, err := server.New(secretClient)
	if err != nil {
		log.Fatalf("could not create server: %s", err)
	}
	address := fmt.Sprintf(":%s", port)

	log.Fatalf("Server stopped: %s", http.ListenAndServe(address, handler))
}
