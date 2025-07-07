package main

import (
	"flag"
	"log/slog"
	"net/http"
	"uptimemonitor/handler"
	"uptimemonitor/router"
	"uptimemonitor/store/sqlite"
)

var (
	dsn  string
	addr string
)

func main() {
	flag.StringVar(&dsn, "dsn", "db.sqlite", "database server name")
	flag.StringVar(&addr, "addr", ":3000", "server address")

	flag.Parse()

	store := sqlite.New(dsn)
	handler := handler.New(store)
	router := router.New(handler)

	server := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	slog.Info("http://localhost:3000")

	server.ListenAndServe()
}
