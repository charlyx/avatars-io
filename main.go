package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/charlyx/avatars.io/secrets"
	server "github.com/charlyx/avatars.io/server"
	lru "github.com/hashicorp/golang-lru"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	cache, err := lru.New(128)
	if err != nil {
		log.Fatalf("could not create cache: %s", err.Error())
	}

	twitterToken, err := secrets.Get("TWITTER_BEARER_TOKEN")
	if err != nil {
		log.Fatalf("could not get twitter token: %s", err.Error())
	}

	handler := server.New(twitterToken, cache)
	address := fmt.Sprintf(":%s", port)

	log.Fatalf("Server stopped: %s", http.ListenAndServe(address, handler))
}
