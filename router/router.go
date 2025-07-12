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

	{
		installedMux := http.NewServeMux()
		installedMux.HandleFunc("GET /login", handler.LoginPage())
		installedMux.HandleFunc("POST /login", handler.LoginForm())

		{
			authenticatedMux := http.NewServeMux()
			authenticatedMux.HandleFunc("GET /{$}", handler.HomePage())

			installedMux.Handle("/",
				handler.UserFromCookie(
					handler.Authenticated(authenticatedMux),
				),
			)
		}

		mux.Handle("/", handler.Installed(installedMux))
	}

	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(static.FS))))

	return mux
}
