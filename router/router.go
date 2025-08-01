package router

import (
	"net/http"
	"uptimemonitor/handler"
	"uptimemonitor/static"
)

func New(handler *handler.Handler, callback func(r *http.ServeMux)) http.Handler {
	r := http.NewServeMux()

	r.HandleFunc("GET /setup", handler.SetupPage())
	r.HandleFunc("POST /setup", handler.SetupForm())
	r.HandleFunc("GET /sponsors", handler.ListSponsors())

	{
		mux := http.NewServeMux()

		{
			loginMux := http.NewServeMux()

			loginMux.HandleFunc("GET /", handler.LoginPage())
			loginMux.HandleFunc("POST /", handler.LoginForm())

			mux.Handle("/login", handler.Guest(loginMux))
		}

		{
			authenticatedMux := http.NewServeMux()

			authenticatedMux.HandleFunc("GET /{$}", handler.HomePage())
			authenticatedMux.HandleFunc("GET /new", handler.CreateMonitorPage())
			authenticatedMux.HandleFunc("GET /monitors", handler.ListMonitors())
			authenticatedMux.HandleFunc("POST /monitors", handler.CreateMonitorForm())
			authenticatedMux.HandleFunc("GET /m/{monitor}", handler.MonitorPage())
			authenticatedMux.HandleFunc("GET /m/{monitor}/i/{incident}", handler.IncidentPage())
			authenticatedMux.HandleFunc("PATCH /monitors/{monitor}", handler.EditMonitorForm())
			authenticatedMux.HandleFunc("DELETE /monitors/{monitor}", handler.DeleteMonitorForm())
			authenticatedMux.HandleFunc("GET /m/{monitor}/edit", handler.EditMonitorPage())
			authenticatedMux.HandleFunc("GET /m/{monitor}/delete", handler.DeleteMonitorPage())
			authenticatedMux.HandleFunc("GET /monitors/{monitor}/checks", handler.ListChecks())
			authenticatedMux.HandleFunc("GET /monitors/{monitor}/stats", handler.MonitorStats())
			authenticatedMux.HandleFunc("GET /monitors/{monitor}/incidents", handler.ListMonitorIncidents())
			authenticatedMux.HandleFunc("GET /incidents", handler.ListIncidents())
			authenticatedMux.HandleFunc("DELETE /incidents/{incident}", handler.DeleteIncident())
			authenticatedMux.HandleFunc("GET /logout", handler.Logout())

			mux.Handle("/", handler.Authenticated(authenticatedMux))
		}

		r.Handle(
			"/",
			handler.UserFromCookie(
				handler.Installed(
					mux,
				),
			),
		)

	}

	r.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(static.FS))))

	callback(r)

	return handler.Recoverer(
		handler.NoCache(r),
	)
}
