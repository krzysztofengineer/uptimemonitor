package app

import (
	"net/http"
	"uptimemonitor/cmd/uptimemonitor/static"
)

func NewRouter(store *Store) *http.ServeMux {
	mux := http.NewServeMux()

	handler := NewHandler(store)

	mux.HandleFunc("/", handler.HomePage())
	mux.HandleFunc("/setup", handler.SetupPage())

	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(static.FS))))

	return mux
}
