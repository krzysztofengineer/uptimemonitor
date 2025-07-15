package handler

import (
	"context"
	"html/template"
	"log/slog"
	"net/http"
	"strconv"
	"sync"
	"time"
	"uptimemonitor"
	"uptimemonitor/html"
	"uptimemonitor/store"
)

type CheckHandler struct {
	Store store.Store
}

func (h *CheckHandler) ListChecks() http.HandlerFunc {
	tmpl := template.Must(template.ParseFS(html.FS, "check.html"))

	type data struct {
		Checks []uptimemonitor.Check
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

		tmpl.ExecuteTemplate(w, "check_list", data{
			Checks: checks,
		})
	}
}

func (h *CheckHandler) RunCheck(ctx context.Context, wg *sync.WaitGroup) error {
	monitors, err := h.Store.ListMonitors(ctx)
	if err != nil {
		return err
	}

	for _, m := range monitors {
		wg.Add(1)

		go func(m uptimemonitor.Monitor) {
			defer wg.Done()

			ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
			defer cancel()

			h.Store.CreateCheck(ctx, uptimemonitor.Check{
				MonitorID: m.ID,
				Monitor:   m,
			})
		}(m)
	}

	return nil
}
