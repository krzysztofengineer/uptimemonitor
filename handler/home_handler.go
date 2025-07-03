package handler

import (
	"net/http"
	"uptimemonitor/store"
)

type HomeHandler struct {
	Store store.Store
}

func (h *HomeHandler) HomePage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/setup", http.StatusSeeOther)
	}
}
