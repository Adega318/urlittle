package main

import (
	"fmt"
	"io"
	"net/http"

	"github.com/Adega318/urlittle/store"
)

type Handler struct {
	st store.Store
}

func NewHandler(st *store.Store) *Handler {
	return &Handler{*st}
}

func (h *Handler) Store(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "failed to read body", http.StatusBadRequest)
		return
	}

	id, err := h.st.Add(r.Context(), string(body))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Fprintf(w, "Stored on: %s", id)
}

func (h *Handler) Retrieve(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	val, err := h.st.Get(r.Context(), id)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	fmt.Fprintf(w, "value: %s", val)
}
