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

	// Converts base64 int to base32
	return id - 76561197960265728, nil
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	opID := steam_go.NewOpenId(r)
	log.Println("Mode", opID.Mode())
	switch opID.Mode() {
	case "":
		http.Redirect(w, r, opID.AuthUrl(), 301)
	case "cancel":
		w.Write([]byte("Authentication cancelled"))
	default:
		steamID, err := opID.ValidateAndGetId()
		if err != nil {
			log.Fatalf("main.go: failed to get steam ID. %s", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		accountID, err := getAccountID(steamID)
		if err != nil {
			log.Fatalf("main.go: failed to get account ID from steam ID. %s", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		fmt.Printf("Account ID: %d", accountID)
		fmt.Fprintln(w, accountID)
	}
}

func matchesHandler(w http.ResponseWriter, r *http.Request) {
	// get accountID from request
	accountID := 2

	// There is no way to query the match with the exact startTime, hence we save the queries immediately.
	matches, err := http.Get(fmt.Sprintf("https://api.opendota.com/api/players/%d/matches?date=7", accountID))
	if err != nil {
		log.Fatalf("main.go: attempt to get recent matches for player with account ID %d failed. %s", accountID, err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	responseBody, err := ioutil.ReadAll(matches.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatalf("main.go: attempt to retrieve the response body from recent matches failed. %s", err.Error())
	}

	if string(responseBody) == "[]" {
		emptyResponseErr := fmt.Sprintf("account ID %d seems to be incorrect as an empty response came back "+
			"attempting to get recent matches.", accountID)
		http.Error(w, emptyResponseErr, http.StatusInternalServerError)
		log.Fatalf("main.go: " + emptyResponseErr)
	} else {
		w.Write([]byte(responseBody))
		// To get date, Date(responseBody[0].start_time)
	}
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("lala")
	enableCors(&w)
	fmt.Fprintln(w, `{
		"userId": 1,
		"id": 1,
		"title": "delectus aut autem",
		"completed": false
	}`)
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

}

func main() {
	fmt.Println("Running...")
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/healthCheck", healthCheckHandler)
	http.HandleFunc("/matches", matchesHandler)
	http.ListenAndServe(":8080", nil)
}
