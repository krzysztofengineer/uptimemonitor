package handler

import (
	"html/template"
	"net/http"
	"uptimemonitor"
	"uptimemonitor/html"
	"uptimemonitor/store"
)

type IncidentHandler struct {
	Store *store.Store
}

func (h *IncidentHandler) ListIncidents() http.HandlerFunc {
	tmpl := template.Must(template.ParseFS(html.FS, "incident.html"))

	type data struct {
		Incidents []uptimemonitor.Incident
	}

	return func(w http.ResponseWriter, r *http.Request) {
		incidents, err := h.Store.ListOpenIncidents(r.Context())
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		tmpl.ExecuteTemplate(w, "incident_list", data{
			Incidents: incidents,
		})
	}
}
