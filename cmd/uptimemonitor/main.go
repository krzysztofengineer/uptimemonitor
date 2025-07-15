package main

import (
	"context"
	"flag"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
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
	flag.StringVar(&dsn, "dsn", "db.sqlite?_journal_mode=WAL&_busy_timeout=5000&_synchronous=FULL&_txlock=immediate", "database server name")
	flag.StringVar(&addr, "addr", ":3000", "server address")

	flag.Parse()

	store := sqlite.New(dsn)
	handler := handler.New(store)
	router := router.New(handler)

	server := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	done := make(chan bool)
	ticker := time.NewTicker(time.Second * 5)

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
				log.Println("ticker")
				handler.RunCheck(context.Background())
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
