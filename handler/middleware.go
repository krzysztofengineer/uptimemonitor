package handler

import (
	"net/http"
	"uptimemonitor/store"
)

type Middleware struct {
	Store store.Store
}

func (m *Middleware) InstalledMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count, err := m.Store.CountUsers(r.Context())
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		if count == 0 && r.Method == http.MethodGet {
			http.Redirect(w, r, "/setup", http.StatusSeeOther)
			return
		}

		if count == 0 {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
