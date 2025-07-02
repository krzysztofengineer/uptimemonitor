package main

import (
	"log/slog"
	"uptimemonitor/app"
	"uptimemonitor/database"
)

func main() {
	db := database.Must(database.New(":memory:"))
	store := app.NewStore(db)

	server := app.NewServer(":3000", store)

	slog.Info("http://localhost:3000")

	server.ListenAndServe()
}
