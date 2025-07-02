package app

import (
	"net/http"
)

type Handler struct {
	Store *Store
}

func NewHandler(store *Store) *Handler {
	return &Handler{Store: store}
}

func (h *Handler) HomePage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/setup", http.StatusSeeOther)
	}
}

func (h *Handler) SetupPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	}
}
