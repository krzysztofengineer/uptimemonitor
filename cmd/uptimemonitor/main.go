package main

import (
	"log/slog"
	"uptimemonitor/app"
	"uptimemonitor/database"
)

func main() {
	db := database.Must(database.New(":memory:"))

	server := app.NewServer(":3000", db)

	slog.Info("http://localhost:3000")

	server.ListenAndServe()
}
