package main

import (
	"context"
	"net/http"

	"github.com/Adega318/urlittle/store"
)

func main() {
	st, _ := store.NewStore(context.Background(), "postgres://urlittle:urlittle@localhost:5432/urlittle", 100)
	defer st.Close()
	h := NewHandler(st)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /", h.Store)
	mux.HandleFunc("GET /{id}", h.Retrieve)

	err := http.ListenAndServe(":9808", mux)
	if err != nil {
		panic(err)
	}
}
