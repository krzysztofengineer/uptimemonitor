package handler

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"uptimemonitor"
	"uptimemonitor/form"
	"uptimemonitor/html"
	"uptimemonitor/store"
)

type MonitorHandler struct {
	Store *store.Store
}

func (h *MonitorHandler) ListMonitors() http.HandlerFunc {
	tmpl := template.Must(template.ParseFS(html.FS, "monitor.html"))

	type data struct {
		Monitors  []uptimemonitor.Monitor
		Skeletons []int
	}

	return func(w http.ResponseWriter, r *http.Request) {
		monitors, err := h.Store.ListMonitors(r.Context())
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		tmpl.ExecuteTemplate(w, "monitor_list", data{
			Monitors:  monitors,
			Skeletons: make([]int, 60),
		})
	}
}

func (h *MonitorHandler) CreateMonitorPage() http.HandlerFunc {
	tmpl := template.Must(template.ParseFS(html.FS, "layout.html", "app.html", "new.html"))

	type data struct {
		Form form.MonitorForm
	}

	return func(w http.ResponseWriter, r *http.Request) {
		tmpl.Execute(w, data{Form: form.MonitorForm{}})
	}
}

func (h *MonitorHandler) CreateMonitorForm() http.HandlerFunc {
	tmpl := template.Must(template.ParseFS(html.FS, "new.html"))

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
			tmpl.ExecuteTemplate(w, "new_form", data{Form: f})
			return
		}

		m, err := h.Store.CreateMonitor(r.Context(), uptimemonitor.Monitor{
			Url: f.Url,
		})
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		w.Header().Set("HX-Redirect", m.URI())
	}
}

func (h *MonitorHandler) ShowMonitor() http.HandlerFunc {
	tmpl := template.Must(template.ParseFS(html.FS, "layout.html", "app.html", "monitor.html"))

	type data struct {
		Monitor   uptimemonitor.Monitor
		Skeletons []int
	}

	return func(w http.ResponseWriter, r *http.Request) {
		m, err := h.Store.GetMonitorByUuid(r.Context(), r.PathValue("monitor"))
		if err != nil || m.ID == 0 {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		tmpl.Execute(w, data{
			Monitor:   m,
			Skeletons: make([]int, 60),
		})
	}
}

func (h *MonitorHandler) MonitorStats() http.HandlerFunc {
	tmpl := template.Must(template.ParseFS(html.FS, "monitor.html"))

	type data struct {
		ID              int64
		AvgResponseTime int64
		Uptime          string
	}

	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("monitor"))
		if err != nil {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		checks, err := h.Store.ListChecks(r.Context(), int64(id), 60)
		if err != nil || id == 0 {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		uptime := float32(0)
		avgResTime := int64(0)

		for _, ch := range checks {
			avgResTime += ch.ResponseTimeMs

			if ch.StatusCode < 400 {
				uptime++
			}
		}

		if len(checks) > 0 {
			avgResTime = avgResTime / int64(len(checks))
			uptime = uptime / float32(len(checks)) * 100.0
		}

		tmpl.ExecuteTemplate(w, "monitor_stats", data{
			ID:              int64(id),
			AvgResponseTime: int64(avgResTime),
			Uptime:          fmt.Sprintf("%.1f", uptime),
		})
	}
}

func (h *MonitorHandler) ListMonitorIncidents() http.HandlerFunc {
	tmpl := template.Must(template.ParseFS(html.FS, "monitor.html"))

	type data struct {
		ID        int64
		Monitor   uptimemonitor.Monitor
		Incidents []uptimemonitor.Incident
	}

	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("monitor"))
		if err != nil {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		m, err := h.Store.GetMonitorByID(r.Context(), id)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		incidents, err := h.Store.ListMonitorIncidents(r.Context(), int64(id))
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		tmpl.ExecuteTemplate(w, "monitor_incident_list", data{
			ID:        int64(id),
			Monitor:   m,
			Incidents: incidents,
		})
	}
}
