package handler

import (
	"uptimemonitor/store"
)

type Handler struct {
	HomeHandler
	SetupHandler
}

func New(store store.Store) *Handler {
	return &Handler{
		HomeHandler:  HomeHandler{Store: store},
		SetupHandler: SetupHandler{Store: store},
	}
}
