package handler

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"
	"uptimemonitor"
	"uptimemonitor/html"
)

func (h *Handler) ListChecks() http.HandlerFunc {
	tmpl := template.Must(template.New("check.html").Funcs(template.FuncMap{
		"sub": func(a, b int) int {
			return a - b
		},
	}).ParseFS(html.FS, "check.html"))

	type data struct {
		Monitor   uptimemonitor.Monitor
		Checks    []uptimemonitor.Check
		Skeletons []int
		MaxTime   int64
		StartTime string
		EndTime   string
	}

	return func(w http.ResponseWriter, r *http.Request) {
		monitorID, err := strconv.Atoi(r.PathValue("monitor"))
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		monitor, err := h.Store.GetMonitorByID(r.Context(), monitorID)
		if err != nil || monitor.ID == 0 {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		checks, err := h.Store.ListChecks(r.Context(), int64(monitorID), 60)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		maxTime := int64(0)
		for _, check := range checks {
			if check.ResponseTimeMs > maxTime {
				maxTime = check.ResponseTimeMs
			}
		}

		err = tmpl.ExecuteTemplate(w, "check_list", data{
			Monitor:   monitor,
			Checks:    checks,
			Skeletons: make([]int, 60),
			MaxTime:   maxTime,
			StartTime: time.Now().Add(-1 * time.Hour).Format("15:04"),
			EndTime:   time.Now().Format("15:04"),
		})

		if err != nil {
			http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
			return
		}
	}
}
