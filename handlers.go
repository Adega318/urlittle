package main

import (
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

	url := "http://" + r.Host + "/" + id

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write([]byte(url))
}

func (h *Handler) Redirect(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	val, err := h.st.Get(r.Context(), id)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	http.Redirect(w, r, val, http.StatusMovedPermanently)
}
