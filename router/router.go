package router

import (
	"net/http"
	"uptimemonitor/handler"
	"uptimemonitor/static"
)

func New(handler *handler.Handler) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /setup", handler.SetupPage())
	mux.HandleFunc("POST /setup", handler.SetupForm())

	installedMux := http.NewServeMux()
	installedMux.HandleFunc("GET /{$}", handler.HomePage())
	installedMux.HandleFunc("GET /login", handler.LoginPage())
	installedMux.HandleFunc("POST /login", handler.LoginForm())

	mux.Handle("/", handler.InstalledMiddleware(installedMux))

	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(static.FS))))

	return mux
}
