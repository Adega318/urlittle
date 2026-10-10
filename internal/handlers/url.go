package handlers

import (
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/Adega318/urlittle/internal/normalize"
	"github.com/Adega318/urlittle/internal/store"
)

const maxBodySize = 8 << 10
const baseUrl = "http://localhost"

func (s *server) store(w http.ResponseWriter, r *http.Request) {
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

	id, err := s.st.Add(r.Context(), value)
	if err != nil {
		slog.ErrorContext(r.Context(), "failed to store URL", "err", err)
		http.Error(w, "failed to store url", http.StatusInternalServerError)
		return
	}

	url := baseUrl + s.addr + "/" + id

	w.Header().Set("Location", url)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	_, _ = io.WriteString(w, url)
}

func (s *server) redirect(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	val, err := s.st.Get(r.Context(), id)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			http.NotFound(w, r)
			return
		}
		slog.ErrorContext(r.Context(), "failed to look up URL",
			"id", id, "err", err)
		http.Error(w, "failed to look up url", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, val, http.StatusTemporaryRedirect)
}
