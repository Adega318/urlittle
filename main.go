package main

import (
	"net/http"
)

func main() {
	h := NewHandler()

	mux := http.NewServeMux()

	mux.HandleFunc("POST /", h.Store)
	mux.HandleFunc("GET /{id}", h.Retrieve)

	err := http.ListenAndServe(":9808", mux)
	if err != nil {
		panic(err)
	}
}
