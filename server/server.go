package server

import (
	"fmt"
	"net/http"

	"github.com/charlyx/avatars.io/secrets"
	"github.com/charlyx/avatars.io/twitter"
	lru "github.com/hashicorp/golang-lru"
)

func Start(port string) error {
	if port == "" {
		port = "8080"
	}

	cache, err := lru.New(128)
	if err != nil {
		return err
	}

	twitterToken, err := secrets.Get("TWITTER_BEARER_TOKEN")
	if err != nil {
		return err
	}

	http.HandleFunc("/", twitter.Handler(twitterToken, cache))

	return http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}
