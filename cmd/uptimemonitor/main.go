package main

import (
	"context"
	"flag"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
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

	var wg sync.WaitGroup
	done := make(chan bool)
	ticker := time.NewTicker(time.Minute)

	go func() {
		slog.Info("http://localhost:3000")

		server.ListenAndServe()
	}()

	go func() {
		for {
			select {
			case <-done:
				slog.Info("done...")
				return
			case <-ticker.C:
				handler.RunCheck(context.Background(), &wg)
			}
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	slog.Info("quitting...")

	done <- true
}
