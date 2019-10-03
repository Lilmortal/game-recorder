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
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/healthCheck", healthCheckHandler)
	http.HandleFunc("/matches", matchesHandler)
	http.ListenAndServe(":8080", nil)

	return nil
}

// TODO:
// Add a middleware for logging
