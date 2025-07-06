package handler

import (
	"html/template"
	"net/http"
	"uptimemonitor/form"
	"uptimemonitor/html"
	"uptimemonitor/store"
)

type LoginHandler struct {
	Store store.Store
}

func (h *LoginHandler) LoginPage() http.HandlerFunc {
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

func (h *LoginHandler) LoginForm() http.HandlerFunc {
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
	}
}
