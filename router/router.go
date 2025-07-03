package router

import (
	"net/http"
	"uptimemonitor/handler"
	"uptimemonitor/static"
)

func New(handler *handler.Handler) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /{$}", handler.HomePage())
	mux.HandleFunc("GET /setup", handler.SetupPage())
	mux.HandleFunc("POST /setup", handler.SetupForm())

	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(static.FS))))

	return mux
}
