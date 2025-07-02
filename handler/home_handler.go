package handler

import (
	"net/http"
	"uptimemonitor"
)

type HomeHandler struct {
	UserStore uptimemonitor.UserStore
}

func (h *HomeHandler) HomePage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/setup", http.StatusSeeOther)
	}
}
