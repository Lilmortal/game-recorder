package main

import (
	"fmt"
	"net/http"
)

// HealthCheckHandler returns a response indicating that this service is still up and running.
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, `{
		"userId": 1,
		"id": 1,
		"title": "delectus aut autem",
		"completed": false
	}`)
}
