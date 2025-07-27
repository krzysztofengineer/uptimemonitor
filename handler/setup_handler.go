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

func (h *Handler) SetupPage() http.HandlerFunc {
	type data struct {
		Form form.SetupForm
	}

	tmpl := template.Must(template.ParseFS(html.FS, "layout.html", "setup.html"))

	return func(w http.ResponseWriter, r *http.Request) {
		count, err := h.Store.CountUsers(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if count > 0 {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		tmpl.Execute(w, data{
			Form: form.SetupForm{},
		})
	}
}

func (h *Handler) SetupForm() http.HandlerFunc {
	type data struct {
		Form form.SetupForm
	}

	tmpl := template.Must(template.ParseFS(html.FS, "setup.html"))

	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()

		f := form.SetupForm{
			Name:     r.PostFormValue("name"),
			Email:    r.PostFormValue("email"),
			Password: r.PostFormValue("password"),
		}

		if !f.Validate() {
			w.WriteHeader(http.StatusBadRequest)
			tmpl.ExecuteTemplate(w, "setup_form", data{
				Form: f,
			})
			return
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(f.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		user, err := h.Store.CreateUser(r.Context(), uptimemonitor.User{
			Name:         f.Name,
			Email:        f.Email,
			PasswordHash: string(hash),
		})
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
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
