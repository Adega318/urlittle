package main

import (
	"fmt"
	"io"
	"net/http"

	lru "github.com/hashicorp/golang-lru/v2"
)

type Handler struct {
	cache *lru.Cache[string, string]
}

func NewHandler(cache *lru.Cache[string, string]) *Handler {
	return &Handler{cache: cache}
}

func (h *Handler) Store(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "failed to read body", http.StatusBadRequest)
		return
	}

	h.cache.Add(id, string(body))

	fmt.Fprintf(w, "Stored on: %s", id)
}
func (h *Handler) Retrieve(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	val, ok := h.cache.Get(id)
	if !ok || val == "" {
		http.NotFound(w, r)
		return
	}

	fmt.Fprintf(w, "value: %s", val)
}

func main() {
	cache, _ := lru.New[string, string](100)
	h := NewHandler(cache)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /{id}", h.Store)
	mux.HandleFunc("GET /{id}", h.Retrieve)

	err := http.ListenAndServe(":9808", mux)
	if err != nil {
		panic(err)
	}
}
