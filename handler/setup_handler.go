package handler

import "net/http"

type SetupHandler struct {
}

func (h *SetupHandler) SetupPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	}
}
