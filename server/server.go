package server

import (
	"log"
	"net/http"
	"os"

	"github.com/charlyx/avatars.io/gravatar"
	"github.com/charlyx/avatars.io/secrets"
	"github.com/charlyx/avatars.io/twitter"
	"github.com/charlyx/avatars.io/usage"
)

func New() (http.Handler, error) {
	secretClient, err := secrets.NewClient(os.Getenv("PROJECT_ID"))
	if err != nil {
		log.Fatalf("could not create secret accessor: %s", err)
	}

	twitterHandlerFunc, err := twitter.NewHandlerFunc(secretClient)
	if err != nil {
		return nil, err
	}

  mux := http.NewServeMux()

	mux.HandleFunc("/twitter", twitterHandlerFunc)
	mux.HandleFunc("/gravatar", gravatar.HandlerFunc)
	mux.HandleFunc("/", usage.HandlerFunc)

	return mux, nil
}
