package handler

import (
	"html/template"
	"net/http"
	"uptimemonitor/html"
)

func (h *Handler) HomePage() http.HandlerFunc {
	tmpl := template.Must(template.ParseFS(html.FS, "layout.html", "app.html", "home.html"))

	return func(w http.ResponseWriter, r *http.Request) {
		count := h.Store.CountMonitors(r.Context())
		if count == 0 {
			http.Redirect(w, r, "/new", http.StatusSeeOther)
			return
		}

		tmpl.Execute(w, nil)
	}
}
