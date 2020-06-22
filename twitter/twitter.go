package twitter

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	lru "github.com/hashicorp/golang-lru"
)

const twitterAPI = "https://api.twitter.com"

type profile struct {
	ImageURL string `json:"profile_image_url"`
}

func getUserProfileImageURL(username, token string) (string, error) {
	userShowURL := fmt.Sprintf("%s/1.1/users/show.json?screen_name=%s", twitterAPI, username)

	req, err := http.NewRequest("GET", userShowURL, nil)
	if err != nil {
		return "", err
	}

	authorization := fmt.Sprintf("Bearer %s", token)
	req.Header.Add("Authorization", authorization)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	userProfile := &profile{}
	if err := json.NewDecoder(resp.Body).Decode(&userProfile); err != nil {
		return "", err
	}

	return userProfile.ImageURL, nil
}

func Handler(token string, cache *lru.Cache) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		username := req.URL.Query().Get("username")

		if username == "" {
			log.Print("no username given")
			http.NotFound(w, req)
			return
		}

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
			newProfileImageURL, err := getUserProfileImageURL(username, token)
			if err != nil {
				http.NotFound(w, req)
				log.Printf("failed getting profile image URL for user %s: %s", username, err.Error())
				return
			}

			if cache != nil {
				cache.Add(username, newProfileImageURL)
				log.Printf("cached URL %s for user %s", newProfileImageURL, username)
			}

			profileImageURL = newProfileImageURL
			log.Printf("got profile image URL %s for user %s", profileImageURL, username)
		}

		http.Redirect(w, req, profileImageURL, http.StatusFound)
	}
}
