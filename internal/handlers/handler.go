package handlers

import (
	"net/http"

	"github.com/Adega318/urlittle/internal/store"
)

type server struct {
	addr string
	st   *store.Store
}

func NewHandler(addr string, st *store.Store) *http.ServeMux {
	s := server{addr, st}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /{$}", s.store)
	mux.HandleFunc("GET /{id}", s.redirect)
	mux.HandleFunc("GET /health", health)

	return mux
}
