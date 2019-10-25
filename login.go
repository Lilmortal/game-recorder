package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/solovev/steam_go"
)

func getAccountID(steamID string) (int, error) {
	id, err := strconv.Atoi(steamID)
	if err != nil {
		return 0, errors.New("failed to convert steamID from a string into an integer")
	}

	// Converts base64 int to base32
	return id - 76561197960265728, nil
}

func createCookie(name string, value string, isHTTPOnly bool) (*http.Cookie, error) {
	// cookieValue := accountIDCookie{
	// 	accountID: accountID,
	// }

	cookie := http.Cookie{
		Name:  name,
		Value: value,
		// Expires in 1 week
		MaxAge:   604800,
		Path:     "/",
		Secure:   false, // TODO: true in test/prod env
		HttpOnly: isHTTPOnly,
		SameSite: http.SameSiteStrictMode,
	}

	return &cookie, nil
}

// LoginHandler handles users attempting to login via Steam.
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Referer())
	r.Host = "localhost:8080"
	r.RequestURI = `/login`

	opID := steam_go.NewOpenId(r)
	switch opID.Mode() {
	case "":
		http.Redirect(w, r, opID.AuthUrl(), 301)
	case "cancel":
		w.Write([]byte("Authentication cancelled"))
	default:
		steamID, err := opID.ValidateAndGetId()
		fmt.Println("Steam ID", steamID)
		if err != nil {
			log.Fatalf("main.go: failed to get steam ID. %s", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		accountID, err := getAccountID(steamID)
		if err != nil {
			log.Fatalf("main.go: failed to get account ID from steam ID. %s", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		cookie, err := createCookie("gameRecorderAccountId", strconv.Itoa(accountID), true)
		if err != nil {
			log.Fatalf("main.go: failed to create cookie. %s", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		http.SetCookie(w, cookie)

		http.Redirect(w, r, "http://localhost:3000", 301)
	}
}