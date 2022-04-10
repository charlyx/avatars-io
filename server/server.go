package server

import (
	"net/http"

	"github.com/charlyx/avatars.io/gravatar"
	"github.com/charlyx/avatars.io/secrets"
	"github.com/charlyx/avatars.io/twitter"
	"github.com/charlyx/avatars.io/usage"
)

func New() (http.Handler, error) {
	twitterHandlerFunc, err := twitter.NewHandlerFunc(secrets.NewEnvClient())
	if err != nil {
		return nil, err
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/twitter", twitterHandlerFunc)
	mux.HandleFunc("/gravatar", gravatar.HandlerFunc)
	mux.HandleFunc("/", usage.HandlerFunc)

	return mux, nil
}
