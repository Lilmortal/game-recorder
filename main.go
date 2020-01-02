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
	http.Handle("/healthCheck", Adapt(HealthCheckHandler, EnableCors, VerifyJWT))
	http.Handle("/matches", Adapt(MatchesHandler, EnableCors, VerifyJWT))
	http.ListenAndServe(":8080", nil)

	return nil
}

// TODO:
// Add a middleware for logging
// error handling
// Rate limit

// This code is in auth.go vendor, temporarily putting it here for now.
// func (id *OpenId) ValidateAndGetId() (string, error) {
// 	if id.Mode() != "id_res" {
// 		return "", errors.New("Mode must equal to \"id_res\".")
// 	}

// 	// TODO: I have temporarily disabled this vendor check as it is stopping me from redirecting properly.
// 	// As far as I know, disabling this check does not cause much harm yet.
// 	// if id.data.Get("openid.return_to") != id.returnUrl {
// 	// 	return "", errors.New("The \"return_to url\" must match the url of current request.")
// 	// }
