package main

import (
	"log/slog"
	"net/http"
	"uptimemonitor/handler"
	"uptimemonitor/router"
	"uptimemonitor/store/sqlite"
)

func main() {
	addr := ":3000"
	store := sqlite.New(":memory:")
	handler := handler.New(store)
	router := router.New(handler)

	server := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	slog.Info("http://localhost:3000")

	server.ListenAndServe()
}
