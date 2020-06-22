package twitter

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/hashicorp/golang-lru/simplelru"
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

func Handler(token string, cache simplelru.LRUCache) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		username := strings.ToLower(req.URL.Query().Get("username"))

		if username == "" {
			log.Print("no username given")
			http.NotFound(w, req)
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

		sizedProfileImageURL := getSizedProfileImageURL(profileImageURL, size)

		http.Redirect(w, req, sizedProfileImageURL, http.StatusFound)
	}
}
