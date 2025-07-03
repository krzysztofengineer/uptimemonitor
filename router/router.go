package router

import (
	"net/http"
	"uptimemonitor/handler"
	"uptimemonitor/static"
)

func New(handler *handler.Handler) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/", handler.HomePage())
	mux.HandleFunc("/setup", handler.SetupPage())

	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(static.FS))))

	return mux
}
