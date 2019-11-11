package main

import (
	"crypto/rand"
	"encoding/base64"
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
	jwtTokenName  = "gameRecordersToken"
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

func createCookie(name string, value string, maxAge int, secure bool, isHTTPOnly bool, isStrict http.SameSite) (*http.Cookie, error) {
	cookie := http.Cookie{
		Name:     name,
		Value:    value,
		MaxAge:   maxAge,
		Path:     "/",
		Secure:   secure,
		HttpOnly: isHTTPOnly,
		SameSite: http.SameSiteStrictMode,
	}

	return &cookie, nil
}

// TODO: put parameters as a struct
func generateJwtToken(cookieName string, value string) (*http.Cookie, error) {
	header := Header{alg: HS256}
	payload := Payload{}

	jwt := Jwt{header: header, payload: payload}
	return createCookie(cookieName, jwt.Build(), secondsInWeek, isNotSecure, isHTTPOnly, http.SameSiteStrictMode)
}

func generateAntiCSRFTokenCookie() (*http.Cookie, error) {
	randomBytes := make([]byte, 32)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, err
	}

	token := base64.URLEncoding.EncodeToString(randomBytes)
	return createCookie("anti-csrf-token", token, secondsInHour, isNotSecure, isNotHTTPOnly, http.SameSiteStrictMode)
}

// LoginHandler handles users attempting to login via Steam.
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	jwtTokenCookie, err := r.Cookie(jwtTokenName)
	if err == nil {
		header := NewHeader()
		jwtToken := NewJWT
	}
	// TODO: Check if JWT exist in request -
	// Verify if JWT header + payload matches with signature
	// Check if JWT access token timeout > current time, if it is, expire
	// Look up PKCE Oauth

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

		cookie, err := generateJwtToken(jwtTokenName, strconv.Itoa(accountID))
		if err != nil {
			log.Fatalf("main.go: failed to create cookie. %s", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		antiCSRFTokenCookie, err := generateAntiCSRFTokenCookie()
		if err != nil {
			log.Fatalf("main.go: failed to create anti CSRF cookie. %s", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		http.SetCookie(w, cookie)
		http.SetCookie(w, antiCSRFTokenCookie)

		http.Redirect(w, r, returnURL, http.StatusMovedPermanently)
	}
}
