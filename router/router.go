package router

import (
	"math/rand"
	"net/http"
	"time"
	"uptimemonitor/handler"
	"uptimemonitor/static"
)

func New(handler *handler.Handler) *http.ServeMux {
	r := http.NewServeMux()

	r.HandleFunc("GET /setup", handler.SetupPage())
	r.HandleFunc("POST /setup", handler.SetupForm())

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
			authenticatedMux.HandleFunc("GET /m/{monitor}", handler.ShowMonitor())
			authenticatedMux.HandleFunc("GET /monitors/{monitor}/checks", handler.ListChecks())
			authenticatedMux.HandleFunc("GET /monitors/{monitor}/stats", handler.MonitorStats())
			authenticatedMux.HandleFunc("GET /incidents", handler.ListIncidents())

			mux.Handle("/", handler.Authenticated(authenticatedMux))
		}

		r.Handle("/", handler.UserFromCookie(
			handler.Installed(mux),
		))
	}

	r.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(static.FS))))

	r.Handle("/random", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		randInt := rand.Intn(100)

		time.Sleep(time.Duration(randInt) * time.Millisecond)

		if randInt < 5 {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}))

	return r
}
