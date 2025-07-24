package handler

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"uptimemonitor"
	"uptimemonitor/form"
	"uptimemonitor/html"
)

func (h *Handler) ListMonitors() http.HandlerFunc {
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

func (h *Handler) CreateMonitorPage() http.HandlerFunc {
	tmpl := template.Must(template.ParseFS(html.FS, "layout.html", "app.html", "new.html"))

	type data struct {
		Form form.MonitorForm
		User uptimemonitor.User
	}

	return func(w http.ResponseWriter, r *http.Request) {
		tmpl.Execute(w, data{
			Form: form.MonitorForm{
				HttpHeaders: `{
	"Content-Type": "application/json"
}`,
				HttpBody: `{}`,
			},
			User: getUserFromRequest(r),
		})
	}
}

func (h *Handler) CreateMonitorForm() http.HandlerFunc {
	tmpl := template.Must(template.ParseFS(html.FS, "new.html"))

	type data struct {
		Form form.MonitorForm
	}

	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()

		f := form.MonitorForm{
			HttpMethod:       r.PostFormValue("http_method"),
			Url:              r.PostFormValue("url"),
			HasCustomHeaders: r.PostFormValue("has_custom_headers") == "on",
			HasCustomBody:    r.PostFormValue("has_custom_body") == "on",
			HttpHeaders:      r.PostFormValue("http_headers"),
			HttpBody:         r.PostFormValue("http_body"),
			HasWebhook:       r.PostFormValue("has_webhook") == "on",
			WebhookMethod:    r.PostFormValue("webhook_method"),
			WebhookUrl:       r.PostFormValue("webhook_url"),
			WebhookHeaders:   r.PostFormValue("webhook_headers"),
			WebhookBody:      r.PostFormValue("webhook_body"),
		}

		if !f.Validate() {
			w.WriteHeader(http.StatusBadRequest)
			tmpl.ExecuteTemplate(w, "new_form", data{Form: f})
			return
		}

		monitor := uptimemonitor.Monitor{
			HttpMethod: f.HttpMethod,
			Url:        f.Url,
		}

		if f.HasCustomHeaders {
			monitor.HttpHeaders = f.HttpHeaders
		}

		if f.HasCustomBody {
			monitor.HttpBody = f.HttpBody
		}

		if f.HasWebhook {
			monitor.WebhookUrl = f.WebhookUrl
			monitor.WebhookMethod = f.WebhookMethod
			monitor.WebhookHeaders = f.WebhookHeaders
			monitor.WebhookBody = f.WebhookBody
		}

		m, err := h.Store.CreateMonitor(r.Context(), monitor)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		w.Header().Set("HX-Redirect", m.URI())
	}
}

func (h *Handler) MonitorPage() http.HandlerFunc {
	tmpl := template.Must(template.ParseFS(html.FS, "layout.html", "app.html", "monitor.html"))

	type data struct {
		Monitor   uptimemonitor.Monitor
		Skeletons []int
		User      uptimemonitor.User
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
			User:      getUserFromRequest(r),
		})
	}
}

func (h *Handler) MonitorStats() http.HandlerFunc {
	tmpl := template.Must(template.ParseFS(html.FS, "monitor.html"))

	type data struct {
		ID              int64
		AvgResponseTime int64
		Uptime          string
		Count           int64
	}

	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("monitor"))
		if err != nil {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		count := h.Store.CountMonitorIncidents(r.Context(), int64(id))

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
			Count:           count,
		})
	}
}

func (h *Handler) ListMonitorIncidents() http.HandlerFunc {
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

func (h *Handler) EditMonitorPage() http.HandlerFunc {
	tmpl := template.Must(template.ParseFS(html.FS, "layout.html", "app.html", "edit.html"))

	type data struct {
		Form    form.MonitorForm
		User    uptimemonitor.User
		Monitor uptimemonitor.Monitor
	}

	return func(w http.ResponseWriter, r *http.Request) {
		m, err := h.Store.GetMonitorByUuid(r.Context(), r.PathValue("monitor"))
		if err != nil {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		f := form.MonitorForm{
			Url:              m.Url,
			HttpMethod:       m.HttpMethod,
			HttpHeaders:      m.HttpHeaders,
			HttpBody:         m.HttpBody,
			HasCustomHeaders: m.HttpHeaders != "",
			HasCustomBody:    m.HttpBody != "",
			HasWebhook:       m.WebhookUrl != "",
			WebhookUrl:       m.WebhookUrl,
			WebhookMethod:    m.WebhookMethod,
			WebhookHeaders:   m.WebhookHeaders,
			WebhookBody:      m.WebhookBody,
		}

		if !f.HasCustomBody {
			f.HttpBody = "{}"
		}

		if !f.HasCustomHeaders {
			f.HttpHeaders = "{}"
		}

		tmpl.Execute(w, data{
			Monitor: m,
			Form:    f,
			User:    getUserFromRequest(r),
		})
	}
}

func (h *Handler) EditMonitorForm() http.HandlerFunc {
	tmpl := template.Must(template.ParseFS(html.FS, "edit.html"))

	type data struct {
		Form    form.MonitorForm
		User    uptimemonitor.User
		Monitor uptimemonitor.Monitor
	}

	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("monitor"))
		if err != nil {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		monitor, err := h.Store.GetMonitorByID(r.Context(), id)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		r.ParseForm()

		f := form.MonitorForm{
			HttpMethod:       r.PostFormValue("http_method"),
			Url:              r.PostFormValue("url"),
			HasCustomHeaders: r.PostFormValue("has_custom_headers") == "on",
			HasCustomBody:    r.PostFormValue("has_custom_body") == "on",
			HttpHeaders:      r.PostFormValue("http_headers"),
			HttpBody:         r.PostFormValue("http_body"),
			HasWebhook:       r.PostFormValue("has_webhook") == "on",
			WebhookMethod:    r.PostFormValue("webhook_method"),
			WebhookUrl:       r.PostFormValue("webhook_url"),
			WebhookHeaders:   r.PostFormValue("webhook_headers"),
			WebhookBody:      r.PostFormValue("webhook_body"),
		}

		if !f.Validate() {
			w.WriteHeader(http.StatusBadRequest)
			tmpl.ExecuteTemplate(w, "edit_form", data{
				Monitor: monitor,
				Form:    f,
				User:    getUserFromRequest(r),
			})
			return
		}

		monitor.Url = f.Url
		monitor.HttpMethod = f.HttpMethod

		if f.HasCustomHeaders {
			monitor.HttpHeaders = f.HttpHeaders
		}

		if f.HasCustomBody {
			monitor.HttpBody = f.HttpBody
		}

		if f.HasWebhook {
			monitor.WebhookUrl = f.WebhookUrl
			monitor.WebhookMethod = f.WebhookMethod
			monitor.WebhookHeaders = f.WebhookHeaders
			monitor.WebhookBody = f.WebhookBody
		}

		err = h.Store.UpdateMonitor(r.Context(), monitor)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		w.Header().Set("HX-Redirect", monitor.URI())
	}
}
