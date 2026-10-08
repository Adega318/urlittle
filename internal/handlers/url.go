package handlers

import (
	"errors"
	"io"
	"net/http"

	"github.com/Adega318/urlittle/internal/normalize"
	"github.com/Adega318/urlittle/internal/store"
)

const maxBodySize = 8 << 10

func (h *Handler) Store(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxBodySize)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		var maxBytesErr *http.MaxBytesError
		if errors.As(err, &maxBytesErr) {
			http.Error(w, "request body too large", http.StatusRequestEntityTooLarge)
			return
		}
		http.Error(w, "failed to read body", http.StatusBadRequest)
		return
	}

	value, err := normalize.URL(string(body))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := h.st.Add(r.Context(), value)
	if err != nil {
		http.Error(w, "failed to store url", http.StatusInternalServerError)
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
		if errors.Is(err, store.ErrNotFound) {
			http.NotFound(w, r)
			return
		}
		http.Error(w, "failed to look up url", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, val, http.StatusMovedPermanently)
}
