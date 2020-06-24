package gravatar

import (
	"crypto/md5"
	"fmt"
	"log"
	"net/http"
	"strings"
)

const gravatarURL = "https://www.gravatar.com/avatar"

func getHash(email string) string {
	normalizedEmail := strings.ToLower(strings.TrimSpace(email))
	md5sum := md5.Sum([]byte(normalizedEmail))
	return fmt.Sprintf("%x", md5sum)
}

func HandlerFunc(w http.ResponseWriter, req *http.Request) {
	email := req.URL.Query().Get("email")

	if email == "" {
		log.Print("no email given")
		http.Error(w, "You must specify email query parameter.", http.StatusBadRequest)
		return
	}

	size := req.URL.Query().Get("size")

	if size == "" {
		s := req.URL.Query().Get("s")

		if s == "" {
			size = "80"
		} else {
			size = s
		}
	}

	hash := getHash(email)
	imageURL := fmt.Sprintf("%s/%s?s=%s", gravatarURL, hash, size)

	log.Printf("computed image URL %s for email %s and size %s", imageURL, email, size)

	http.Redirect(w, req, imageURL, http.StatusFound)
}
