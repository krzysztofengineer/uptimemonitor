package app

import (
	"net/http"
	"uptimemonitor/handler"
	"uptimemonitor/static"
)

type Router struct {
	*http.ServeMux
}

func NewRouter(handler *handler.Handler) *Router {
	mux := http.NewServeMux()

	mux.HandleFunc("/", handler.HomePage())
	mux.HandleFunc("/setup", handler.SetupPage())

	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(static.FS))))

	return &Router{
		ServeMux: mux,
	}
}
