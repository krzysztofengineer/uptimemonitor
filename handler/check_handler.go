package handler

import (
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"strconv"
	"uptimemonitor"
	"uptimemonitor/html"
	"uptimemonitor/store"
)

type CheckHandler struct {
	Store store.Store
}

func (h *CheckHandler) ListChecks() http.HandlerFunc {
	tmpl := template.Must(template.New("check.html").Funcs(template.FuncMap{
		"sub": func(a, b int) int {
			return a - b
		},
	}).ParseFS(html.FS, "check.html"))

	type data struct {
		Monitor   uptimemonitor.Monitor
		Checks    []uptimemonitor.Check
		Skeletons []int
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
			slog.Error("list checks error", "err", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		err = tmpl.ExecuteTemplate(w, "check_list", data{
			Monitor:   monitor,
			Checks:    checks,
			Skeletons: make([]int, 60),
		})

		if err != nil {
			slog.Error("template execution error", "err", err)
			http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
			return
		}
	}
}
