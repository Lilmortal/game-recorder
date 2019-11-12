package main

import (
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
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		h.ServeHTTP(w, r)
	})
}
