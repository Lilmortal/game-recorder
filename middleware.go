package main

import (
	"fmt"
	"log"
	"net/http"
)

// Adapter allows handlers to inject customizable middleware functions.
type Adapter func(http.Handler) http.Handler

// Adapt is a helper function which composes a list of middlewares.
func Adapt(handleFunc func(http.ResponseWriter, *http.Request), adapters ...Adapter) http.Handler {
	var handler http.Handler
	handler = http.HandlerFunc(handleFunc)
	for _, adapter := range adapters {
		handler = adapter(handler)
	}

	return handler
}

// EnableCors is a middleware that enables CORS to all incoming requests.
func EnableCors(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: Remember to change this
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		h.ServeHTTP(w, r)
	})
}

// VerifyJWT ensures all incoming requests after logging in is a verified user via JWT, as well as ensuring
// that it is not tampered by CSRF attacks.
func VerifyJWT(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer h.ServeHTTP(w, r)

		fmt.Println(r.Cookies())
		antiCSRFTokenCookie, err := r.Cookie(AntiCSRFTokenKey)
		if err != nil {
			log.Fatal("Anti CSRF token cookie is missing.")
		}

		antiCSRFTokenHeader := r.Header.Get(xCSRFTokenHeader)

		jwtTokenCookie, err := r.Cookie(JWTTokenKey)

		if err != nil {
			log.Fatal("JWT token cookie is missing.")
		}

		jwt := GetJWT(jwtTokenCookie.Value)
		if jwt.Verify() && isCookieTampered(antiCSRFTokenCookie.Value, antiCSRFTokenHeader) {
			//&& jwt.payload.exp > time.Now().String()
			// TODO: PKCE OAuth
			// get new token
		} else {
			// log.Fatalf("main.go: failed to get jwt token cookie. %s", err.Error())
			// http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Fatal("Failed to verify")
		}
	})
}
