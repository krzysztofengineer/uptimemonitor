package handler

import (
	"html/template"
	"net/http"
	"uptimemonitor/html"
	"uptimemonitor/store"
)

type HomeHandler struct {
	Store store.Store
}

func (h *HomeHandler) HomePage() http.HandlerFunc {
	tmpl := template.Must(template.ParseFS(html.FS, "layout.html", "app.html", "home.html"))

	return func(w http.ResponseWriter, r *http.Request) {
		tmpl.Execute(w, nil)
	}
}
