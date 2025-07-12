package handler

import (
	"html/template"
	"net/http"
	"uptimemonitor"
	"uptimemonitor/html"
	"uptimemonitor/store"
)

type MonitorHandler struct {
	Store store.Store
}

func (h *MonitorHandler) ListMonitors() http.HandlerFunc {
	tmpl := template.Must(template.ParseFS(html.FS, "monitor.html"))

	type data struct {
		Monitors []uptimemonitor.Monitor
	}

	return func(w http.ResponseWriter, r *http.Request) {
		monitors, err := h.Store.ListMonitors(r.Context())
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		tmpl.ExecuteTemplate(w, "monitor_list", data{
			Monitors: monitors,
		})
	}
}
