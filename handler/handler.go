package handler

import (
	"uptimemonitor/store"
)

type Handler struct {
	Middleware
	HomeHandler
	SetupHandler
	LoginHandler
}

func New(store store.Store) *Handler {
	return &Handler{
		Middleware:   Middleware{Store: store},
		HomeHandler:  HomeHandler{Store: store},
		SetupHandler: SetupHandler{Store: store},
		LoginHandler: LoginHandler{Store: store},
	}
}
