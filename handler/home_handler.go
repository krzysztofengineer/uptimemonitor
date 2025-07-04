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
		count, err := h.Store.CountUsers(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if count == 0 {
			http.Redirect(w, r, "/setup", http.StatusSeeOther)
		}
	}
}
