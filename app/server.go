package app

import (
	"net/http"
)

func NewServer(addr string, store *Store) *http.Server {
	return &http.Server{
		Addr:    addr,
		Handler: NewRouter(NewHandler(store)),
	}
}
