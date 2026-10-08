package handlers

import "github.com/Adega318/urlittle/internal/store"

type Handler struct {
	st *store.Store
}

func NewHandler(st *store.Store) *Handler {
	return &Handler{st: st}
}
