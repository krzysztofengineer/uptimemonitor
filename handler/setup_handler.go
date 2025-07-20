package handler

import (
	"html/template"
	"net/http"
	"uptimemonitor"
	"uptimemonitor/form"
	"uptimemonitor/html"
	"uptimemonitor/store"

	"golang.org/x/crypto/bcrypt"
)

type SetupHandler struct {
	Store *store.Store
}

func (h *SetupHandler) SetupPage() http.HandlerFunc {
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

func (h *SetupHandler) SetupForm() http.HandlerFunc {
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

		h.Store.CreateUser(r.Context(), uptimemonitor.User{
			Name:         f.Name,
			Email:        f.Email,
			PasswordHash: string(hash),
		})

		w.Header().Set("HX-Redirect", "/")
	}
}
