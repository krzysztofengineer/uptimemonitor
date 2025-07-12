package handler

import (
	"net/http"
	"uptimemonitor/store"
)

type MonitorHandler struct {
	Store store.Store
}

func (h *MonitorHandler) ListMonitors() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
