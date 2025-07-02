package app

import (
	"net/http"
)

type Server struct {
	*http.Server
}

func NewServer(addr string, router *Router) *Server {
	return &Server{
		Server: &http.Server{
			Addr:    addr,
			Handler: router,
		},
	}
}
