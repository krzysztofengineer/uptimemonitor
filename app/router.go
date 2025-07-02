package app

import (
	"database/sql"
	"net/http"
	"uptimemonitor/cmd/uptimemonitor/static"
)

func NewRouter(db *sql.DB) *http.ServeMux {
	mux := http.NewServeMux()

	handler := NewHandler(db)

	mux.HandleFunc("/", handler.HomePage())
	mux.HandleFunc("/setup", handler.SetupPage())

	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(static.FS))))

	return mux
}
