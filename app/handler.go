package app

import (
	"uptimemonitor"
	"uptimemonitor/handler"
)

type Handler struct {
	handler.HomeHandler
	handler.SetupHandler
}

func NewHandler(store uptimemonitor.Store) *Handler {
	return &Handler{
		HomeHandler:  handler.HomeHandler{UserStore: store},
		SetupHandler: handler.SetupHandler{},
	}
}
