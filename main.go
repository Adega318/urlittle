package main

import (
	"context"
	"net/http"

	"github.com/Adega318/urlittle/store"
)

func main() {
	config, err := LoadConfig()
	if err != nil {
		panic(err)
	}

	st, _ := store.NewStore(context.Background(), config.Store.DBURL, config.Store.CacheSize)
	defer st.Close()
	h := NewHandler(st)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /{$}", h.Store)
	mux.HandleFunc("/{id}", h.Redirect)

	err = http.ListenAndServe(":"+config.Port, mux)
	if err != nil {
		panic(err)
	}
}
