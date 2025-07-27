package handler

import (
	"html/template"
	"net/http"
	"time"
	"uptimemonitor"
	"uptimemonitor/form"
	"uptimemonitor/html"

	"golang.org/x/crypto/bcrypt"
)

func (h *Handler) LoginPage() http.HandlerFunc {
	type data struct {
		Form form.LoginForm
	}

	tmpl := template.Must(template.ParseFS(html.FS, "layout.html", "login.html"))

	return func(w http.ResponseWriter, r *http.Request) {
		tmpl.Execute(w, data{
			Form: form.LoginForm{},
		})
	}
}

func (h *Handler) LoginForm() http.HandlerFunc {
	type data struct {
		Form form.LoginForm
	}

	tmpl := template.Must(template.ParseFS(html.FS, "login.html"))

	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()

		f := form.LoginForm{
			Email:    r.PostFormValue("email"),
			Password: r.PostFormValue("password"),
		}

		if !f.Validate() {
			w.WriteHeader(http.StatusBadRequest)
			tmpl.ExecuteTemplate(w, "login_form", data{
				Form: f,
			})
			return
		}

		user, err := h.Store.GetUserByEmail(r.Context(), f.Email)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			f.Errors["Email"] = "The credentials do not match our records"
			tmpl.ExecuteTemplate(w, "login_form", data{
				Form: f,
			})
			return
		}

		if err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(f.Password)); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			f.Errors["Email"] = "The credentials do not match our records"
			tmpl.ExecuteTemplate(w, "login_form", data{
				Form: f,
			})
			return
		}

		session, err := h.Store.CreateSession(r.Context(), uptimemonitor.Session{
			UserID:    user.ID,
			ExpiresAt: time.Now().Add(time.Hour * 24 * 30),
			User:      user,
		})
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			f.Errors["Email"] = "Something went wrong, try again later"
			tmpl.ExecuteTemplate(w, "login_form", data{
				Form: f,
			})
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "session",
			Value:    session.Uuid,
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
			Secure:   h.Secure,
			Expires:  session.ExpiresAt,
		})

		w.Header().Set("HX-Redirect", "/")
	}
}
