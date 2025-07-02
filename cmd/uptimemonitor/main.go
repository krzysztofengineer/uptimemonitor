package main

import (
	"log/slog"
	"uptimemonitor/app"
	"uptimemonitor/sqlite"
)

func main() {
	store := sqlite.New(":memory:")
	handler := app.NewHandler(store)
	router := app.NewRouter(handler)
	server := app.NewServer(":3000", router)

	slog.Info("http://localhost:3000")

	server.ListenAndServe()
}
