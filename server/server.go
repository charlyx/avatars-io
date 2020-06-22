package server

import (
	"fmt"
	"net/http"

	"github.com/charlyx/avatars.io/secrets"
	"github.com/charlyx/avatars.io/twitter"
	lru "github.com/hashicorp/golang-lru"
)

const usage = `<html><head><title>Not found</title></head><body>
<h1>Not found</h1>

<p>
Give a username and get an avatar in return: <a href="https://avatars.charlyx.dev/twitter?username=charlyx">https://avatars.charlyx.dev/twitter?username=charlyx</a>
</p>

<p>
You can ask for variant sizings such as "bigger", "mini" and "original" (default size being "normal").
</p>

<p>
Example: <a href="https://avatars.charlyx.dev/twitter?username=charlyx&size=bigger">https://avatars.charlyx.dev/twitter?username=charlyx&size=bigger</a>
</p>
</body>
</html>`

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, usage)
}

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

	http.HandleFunc("/twitter", twitter.Handler(twitterToken, cache))
	http.HandleFunc("/", defaultHandler)

	return http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}
