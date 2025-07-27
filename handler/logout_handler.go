package handler

import (
	"net/http"
	"time"
	"uptimemonitor"
)

func (h *Handler) Logout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, ok := r.Context().Value(sessionContextKey).(uptimemonitor.Session)
		if !ok {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusServiceUnavailable)
			return
		}

		err := h.Store.RemoveSessionByID(r.Context(), session.ID)
		if err != nil {
			// todo: log
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "session",
			Value:    "",
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
			Secure:   h.Secure,
			Expires:  time.Now().Add(time.Hour * 24 * 30 * -1),
		})

		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}
