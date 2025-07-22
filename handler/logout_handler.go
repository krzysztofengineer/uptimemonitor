package handler

import (
	"net/http"
	"uptimemonitor"
)

func (h *Handler) Logout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// http.SetCookie(w, &http.Cookie{
		// 	Name:     "session",
		// 	Value:    "",
		// 	HttpOnly: true,
		// 	SameSite: http.SameSiteLaxMode,
		// 	Secure:   false, // todo
		// 	Expires:  time.Now().Add(time.Hour * 24 * 30 * -1),
		// })

		session, ok := r.Context().Value(sessionContextKey).(uptimemonitor.Session)
		if !ok {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		err := h.Store.RemoveSessionByID(r.Context(), session.ID)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}
