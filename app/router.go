package app

import (
	"net/http"
	"uptimemonitor/static"
)

func NewRouter(handler *Handler) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/", handler.HomePage())
	mux.HandleFunc("/setup", handler.SetupPage())

	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(static.FS))))

	return mux
}
