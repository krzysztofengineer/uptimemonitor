package handler

import (
	"html/template"
	"net/http"
	"uptimemonitor/html"
	"uptimemonitor/store"
)

type SetupHandler struct {
	Store store.Store
}

func (h *SetupHandler) SetupPage() http.HandlerFunc {
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

		tmpl.Execute(w, nil)
	}
}

func (h *SetupHandler) SetupForm() http.HandlerFunc {
	tmpl := template.Must(template.ParseFS(html.FS, "setup.html"))

	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)

		tmpl.ExecuteTemplate(w, "setup_form", nil)
	}
}
