package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/solovev/steam_go"
)

const (
	secondsInHour = 3600
	secondsInWeek = 604800
	isSecure      = true
	isNotSecure   = false
	isHTTPOnly    = true
	isNotHTTPOnly = false
)

func convertBase64ToBase32(number int) int {
	return number - 76561197960265728
}

func getAccountID(steamID string) (int, error) {
	id, err := strconv.Atoi(steamID)
	if err != nil {
		return 0, errors.New("failed to convert steamID from a string into an integer")
	}

	return convertBase64ToBase32(id), nil
}

// LoginHandler handles users attempting to login via Steam.
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	jwtTokenCookie, err := r.Cookie(JWTTokenKey)
	if err == nil {
		jwt := GetJWT(jwtTokenCookie.Value)
		if jwt.Verify() {
			//&& jwt.payload.exp > time.Now().String()
			// TODO: verify via CSRF token
			// TODO: PKCE OAuth
			// get new token
		} else {
			// log.Fatalf("main.go: failed to get jwt token cookie. %s", err.Error())
			// http.Error(w, err.Error(), http.StatusInternalServerError)
			fmt.Println("Failed to verify")
		}
	}

	returnURL := ""
	if len(r.URL.Query()["referer"]) > 0 {
		returnURL = r.URL.Query()["referer"][0]
	}

	r.Host = "localhost:8080"
	r.RequestURI = `/login?referer=` + r.Referer()

	opID := steam_go.NewOpenId(r)

	switch opID.Mode() {
	case "":
		http.Redirect(w, r, opID.AuthUrl(), http.StatusMovedPermanently)
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

		cookie, err := GenerateJWTTokenCookie(JWTTokenKey, strconv.Itoa(accountID))
		if err != nil {
			log.Fatalf("main.go: failed to create cookie. %s", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		antiCSRFTokenCookie, err := GenerateAntiCSRFTokenCookie(AntiCSRFTokenKey)
		if err != nil {
			log.Fatalf("main.go: failed to create anti CSRF cookie. %s", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		http.SetCookie(w, cookie)
		http.SetCookie(w, antiCSRFTokenCookie)

		http.Redirect(w, r, returnURL, http.StatusMovedPermanently)
	}
}
