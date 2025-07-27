package handler

import (
	"net/http"
	"uptimemonitor"
	"uptimemonitor/service"
	"uptimemonitor/store"
)

type Handler struct {
	Store  *store.Store
	Secure bool
}

func New(store *store.Store, service *service.Service, secure bool) *Handler {
	return &Handler{
		Store:  store,
		Secure: secure,
	}
}

func getUserFromRequest(r *http.Request) uptimemonitor.User {
	user, ok := r.Context().Value(userContextKey).(uptimemonitor.User)
	if !ok {
		return uptimemonitor.User{}
	}

	return user
}
