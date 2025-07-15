package handler

import (
	"context"
	"html/template"
	"log"
	"log/slog"
	"net/http"
	"strconv"
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

func (h *CheckHandler) RunCheck(ctx context.Context) error {
	monitors, err := h.Store.ListMonitors(ctx)
	if err != nil {
		return err
	}

	log.Printf("running check: %d", len(monitors))

	for _, m := range monitors {
		go func(mon uptimemonitor.Monitor) {
			c, cancel := context.WithTimeout(context.Background(), time.Minute)
			defer cancel()

			log.Printf("CHECK #%d", m.ID)

			check, err := h.Store.CreateCheck(c, uptimemonitor.Check{
				MonitorID: mon.ID,
				Monitor:   mon,
			})
			if err != nil {
				log.Printf("err: #%v", err)
				return
			}

			log.Printf("CHECK FINISHED WITH ID: #%d", check.ID)
		}(m)
	}

	return nil
}
