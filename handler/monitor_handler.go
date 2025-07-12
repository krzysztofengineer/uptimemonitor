package handler

import (
	"html/template"
	"net/http"
	"uptimemonitor"
	"uptimemonitor/form"
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

func (h *MonitorHandler) CreateMonitor() http.HandlerFunc {
	tmpl := template.Must(template.ParseFS(html.FS, "monitor.html"))

	type data struct {
		Form form.MonitorForm
	}

	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()

		f := form.MonitorForm{
			Url: r.PostFormValue("url"),
		}

		if !f.Validate() {
			w.WriteHeader(http.StatusBadRequest)
			tmpl.ExecuteTemplate(w, "monitor_form", data{Form: f})
			return
		}
	}
}
