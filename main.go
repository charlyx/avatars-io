package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/charlyx/avatars.io/server"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	address := fmt.Sprintf(":%s", port)

	handler, err := server.New()
	if err != nil {
		log.Fatalf("could not create server: %s", err)
	}

	log.Fatalf("Server stopped: %s", http.ListenAndServe(address, handler))
}
