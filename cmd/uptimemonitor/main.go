package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"uptimemonitor/handler"
	"uptimemonitor/router"
	"uptimemonitor/service"
	"uptimemonitor/store"
)

var (
	dsn    string
	addr   string
	secure bool
)

func main() {
	flag.StringVar(&dsn, "dsn", "uptimemonitor.sqlite?_journal_mode=WAL&_busy_timeout=5000&_synchronous=FULL&_txlock=immediate", "database server name")
	flag.StringVar(&addr, "addr", ":3000", "server address")
	flag.BoolVar(&secure, "secure", true, "use https")

	flag.Parse()

	store := store.New(dsn)
	service := service.New(store)
	handler := handler.New(store, service, secure)
	router := router.New(handler, func(r *http.ServeMux) {})

	server := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	done := make(chan bool)
	ticker := time.NewTicker(time.Second * 60)

	go func() {
		slog.Info(fmt.Sprintf("http://localhost%s", addr), "secure", secure)

		server.ListenAndServe()
	}()

	checkCh := service.StartCheck()

	go func() {
		service.RunCheck(context.Background(), checkCh)
	}()

	// todo: move to service
	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				service.RunCheck(context.Background(), checkCh)
			}
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	time.Sleep(3 * time.Second)

	done <- true
}
