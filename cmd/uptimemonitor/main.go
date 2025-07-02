package main

import (
	"log/slog"
	"net/http"
	"uptimemonitor/static"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`
		<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>Uptime Monitor</title>
			<link rel="stylesheet" href="/static/css/main.css">
			<script src="/static/js/main.js"></script>
		</head>
		<body>
			<h1>Uptime Monitor</h1>
		</body>
		</html>
		`))
	})

	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(static.FS))))

	server := &http.Server{
		Addr:    ":3000",
		Handler: mux,
	}

	slog.Info("http://localhost:3000")

	server.ListenAndServe()
}
