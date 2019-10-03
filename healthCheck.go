package main

import (
	"fmt"
	"net/http"
)

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
