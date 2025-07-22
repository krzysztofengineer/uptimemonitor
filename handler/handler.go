package handler

import (
	"uptimemonitor/service"
	"uptimemonitor/store"
)

type Handler struct {
	Store *store.Store
}

func New(store *store.Store, service *service.Service) *Handler {
	return &Handler{
		Store: store,
	}
}
