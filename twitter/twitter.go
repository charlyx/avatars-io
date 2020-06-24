package twitter

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/charlyx/avatars.io/secrets"
	lru "github.com/hashicorp/golang-lru"
	"github.com/hashicorp/golang-lru/simplelru"
)

const DefaultImageURL = "https://abs.twimg.com/sticky/default_profile_images/default_profile_normal.png"
const ShowURL = "https://api.twitter.com/1.1/users/show.json"

type profile struct {
	ImageURL string `json:"profile_image_url"`
}

func getUserProfileImageURL(username, token string) string {
	userShowURL := fmt.Sprintf("%s?screen_name=%s", ShowURL, username)

	req, err := http.NewRequest("GET", userShowURL, nil)
	if err != nil {
		log.Printf("got an error creating GET request for %s", userShowURL)
		return DefaultImageURL
	}

	authorization := fmt.Sprintf("Bearer %s", token)
	req.Header.Add("Authorization", authorization)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("got an error requesting profile image URL for user %s", username)
		return DefaultImageURL
	}
	defer resp.Body.Close()

	userProfile := &profile{}
	if err := json.NewDecoder(resp.Body).Decode(&userProfile); err != nil {
		log.Printf("got an error decoding user %s profile", username)
		return DefaultImageURL
	}

	if userProfile.ImageURL == "" {
		return DefaultImageURL
	}

	return userProfile.ImageURL
}

func getNormalizedSize(size string) string {
	if size != "bigger" && size != "mini" && size != "normal" && size != "original" {
		return "normal"
	}
	return strings.ToLower(size)
}

func getSizedProfileImageURL(imageURL, size string) string {
	size = getNormalizedSize(size)

	if size == "normal" {
		return imageURL
	}

	if size == "original" {
		return strings.Replace(imageURL, "_normal", "", 1)
	}

	return strings.Replace(imageURL, "normal", size, 1)
}

func NewHandlerFunc(secret secrets.SecretAccessor) (http.HandlerFunc, error) {
	cache, err := lru.New(128)
	if err != nil {
		return nil, fmt.Errorf("could not create cache: %s", err.Error())
	}

	token, err := secret.Get("TWITTER_BEARER_TOKEN")
	if err != nil {
		return nil, fmt.Errorf("could not get twitter token: %s", err.Error())
	}

	return handlerFunc(token, cache), nil
}

func handlerFunc(token string, cache simplelru.LRUCache) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		username := strings.ToLower(req.URL.Query().Get("username"))

		if username == "" {
			log.Print("no username given")
			http.Error(w, "You must specify username query parameter.", http.StatusBadRequest)
			return
		}

		size := getNormalizedSize(req.URL.Query().Get("size"))

		profileImageURL := ""

		if cache != nil {
			if cachedProfileImageURL, ok := cache.Get(username); ok {
				profileImageURL = cachedProfileImageURL.(string)

				log.Printf("looked up cache for user %s, got URL %s", username, profileImageURL)
			} else {
				log.Printf("got error retrieving cache for user %s", username)
			}
		}

		if profileImageURL == "" {
			profileImageURL = getUserProfileImageURL(username, token)

			log.Printf("got profile image URL %s for user %s", profileImageURL, username)

			if cache != nil {
				cache.Add(username, profileImageURL)

				log.Printf("cached URL %s for user %s", profileImageURL, username)
			}
		}

		sizedProfileImageURL := getSizedProfileImageURL(profileImageURL, size)

		http.Redirect(w, req, sizedProfileImageURL, http.StatusFound)
	}
}
