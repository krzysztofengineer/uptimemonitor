package app

import (
	"uptimemonitor/handler"
)

type Handler struct {
	handler.HomeHandler
	handler.SetupHandler
}

func NewHandler(store Store) *Handler {
	return &Handler{
		HomeHandler:  handler.HomeHandler{UserStore: store},
		SetupHandler: handler.SetupHandler{},
	}
}
