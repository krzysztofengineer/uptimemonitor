package handler

import (
	"net/http"
	"uptimemonitor/store"
)

type LoginHandler struct {
	Store store.Store
}

func (h *LoginHandler) LoginPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func (h *LoginHandler) LoginForm() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
