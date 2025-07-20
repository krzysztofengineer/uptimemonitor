package handler

import (
	"uptimemonitor/service"
	"uptimemonitor/store"
)

type Handler struct {
	Middleware
	HomeHandler
	SetupHandler
	LoginHandler
	MonitorHandler
	CheckHandler
	IncidentHandler
}

func New(store *store.Store, service *service.Service) *Handler {
	return &Handler{
		Middleware:      Middleware{Store: store},
		HomeHandler:     HomeHandler{Store: store},
		SetupHandler:    SetupHandler{Store: store},
		LoginHandler:    LoginHandler{Store: store},
		MonitorHandler:  MonitorHandler{Store: store},
		CheckHandler:    CheckHandler{Store: store},
		IncidentHandler: IncidentHandler{Store: store},
	}
}
