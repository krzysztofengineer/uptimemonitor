package app

import (
	"database/sql"
	"net/http"
)

func NewServer(addr string, db *sql.DB) *http.Server {
	return &http.Server{
		Addr:    addr,
		Handler: NewRouter(db),
	}
}
