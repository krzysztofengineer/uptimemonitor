package handler

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"uptimemonitor"
	"uptimemonitor/html"
)

func (h *Handler) ListIncidents() http.HandlerFunc {
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

func (h *Handler) DeleteIncident() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("incident"))
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		i, err := h.Store.GetIncidentByID(r.Context(), int64(id))
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		h.Store.DeleteIncident(r.Context(), int64(id))

		w.Header().Set("HX-Redirect", fmt.Sprintf("/m/%s", i.Monitor.Uuid))
	}
}

func (h *Handler) IncidentPage() http.HandlerFunc {
	tmpl := template.Must(template.ParseFS(html.FS, "layout.html", "app.html", "incident.html"))

	type data struct {
		User     uptimemonitor.User
		Incident uptimemonitor.Incident
		Monitor  uptimemonitor.Monitor
	}

	return func(w http.ResponseWriter, r *http.Request) {
		uuid := r.PathValue("incident")
		incident, err := h.Store.GetIncidentByUuid(r.Context(), uuid)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		if incident.Monitor.Uuid != r.PathValue("monitor") {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		tmpl.Execute(w, data{
			User:     getUserFromRequest(r),
			Incident: incident,
			Monitor:  incident.Monitor,
		})
	}
}
