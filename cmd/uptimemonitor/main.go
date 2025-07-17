package main

import (
	"context"
	"flag"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"uptimemonitor/handler"
	"uptimemonitor/router"
	"uptimemonitor/service"
	"uptimemonitor/store/sqlite"
)

var (
	dsn  string
	addr string
)

func main() {
	flag.StringVar(&dsn, "dsn", "db.sqlite?_journal_mode=WAL&_busy_timeout=5000&_synchronous=FULL&_txlock=immediate", "database server name")
	flag.StringVar(&addr, "addr", ":3000", "server address")

	flag.Parse()

	store := sqlite.New(dsn)
	handler := handler.New(store)
	service := service.New(store)
	router := router.New(handler)

	server := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	done := make(chan bool)
	ticker := time.NewTicker(time.Minute)

	go func() {
		slog.Info("http://localhost:3000")

		server.ListenAndServe()
	}()

	checkCh := service.StartCheck()

	go func() {
		service.RunChecks(context.Background(), checkCh)
	}()

	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				service.RunChecks(context.Background(), checkCh)
			}
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	slog.Info("quitting...")

	done <- true

	// todo add maximum time to wait
}
