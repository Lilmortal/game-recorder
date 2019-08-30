package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/solovev/steam_go"
)

func substr(input string, start int, length int) string {
	asRunes := []rune(input)

	if start >= len(asRunes) {
		return ""
	}

	if start+length > len(asRunes) {
		length = len(asRunes) - start
	}

	return string(asRunes[start : start+length])
}

func getAccountID(steamID string) (int, error) {
	id, err := strconv.Atoi(steamID)
	if err != nil {
		return 0, errors.New("failed to convert steamID from a string into an integer")
	}
	return id - 765611979602657280, nil
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	opID := steam_go.NewOpenId(r)
	switch opID.Mode() {
	case "":
		http.Redirect(w, r, opID.AuthUrl(), 301)
	case "cancel":
		w.Write([]byte("Authentication cancelled"))
	default:
		steamID, err := opID.ValidateAndGetId()
		if err != nil {
			log.Printf("main.go: failed to get steam ID. \n %s", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		accountID, err := getAccountID(steamID)
		if err != nil {
			log.Printf("main.go: failed to get account ID from steam ID. \n %s", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		matches, err := http.Get(fmt.Sprintf("https://api.opendota.com/api/players/%d/recentMatches", accountID))
		if err != nil {
			log.Printf("main.go: attempt to get recent matches for player with account ID %d failed. \n %s", accountID, err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		responseBody, err := ioutil.ReadAll(matches.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Printf("main.go: attempt to retrieve the response body from recent matches failed. \n %s", err.Error())
		}

		if string(responseBody) == "[]" {
			emptyResponseErr := fmt.Sprintf("account ID %d seems to be incorrect as an empty response came back "+
				"attempting to get recent matches.", accountID)
			http.Error(w, emptyResponseErr, http.StatusInternalServerError)
			log.Printf("main.go: " + emptyResponseErr)
		} else {
			w.Write([]byte(responseBody))
			// To get date, Date(responseBody[0].start_time)
		}
	}
}

func main() {
	http.HandleFunc("/login", loginHandler)
	http.ListenAndServe(":8080", nil)
}
