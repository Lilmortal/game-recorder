package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	fmt.Println("Running...")

	if err := run(); err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}
}

func run() error {
	http.Handle("/login", Adapt(LoginHandler, EnableCors))
	http.Handle("/healthCheck", Adapt(HealthCheckHandler, EnableCors))
	http.Handle("/matches", Adapt(MatchesHandler, EnableCors))
	http.ListenAndServe(":8080", nil)

	return nil
}

// TODO:
// Add a middleware for logging
// error handling
// Verifying all cookies by checking if CSRF token in HTTP header is same as cookie
