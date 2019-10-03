package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

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
