package handler

import (
	"context"
	"net/http"
	"uptimemonitor"
	"uptimemonitor/store"
)

type Middleware struct {
	Store store.Store
}

func (m *Middleware) Installed(next http.Handler) http.Handler {
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

func (m *Middleware) UserFromCookie(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("session")
		if err != nil || c.Value == "" {
			next.ServeHTTP(w, r)
			return
		}

		session, err := m.Store.GetSessionByUuid(r.Context(), c.Value)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		ctx := context.WithValue(r.Context(), "session", session)
		ctx = context.WithValue(ctx, "user", session.User)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *Middleware) Authenticated(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		value := r.Context().Value("session")
		if value == nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		_, ok := value.(uptimemonitor.Session)
		if !ok {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		next.ServeHTTP(w, r)
	})
}
