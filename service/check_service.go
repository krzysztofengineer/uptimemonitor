package service

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
	"uptimemonitor"
	"uptimemonitor/store"
)

type CheckService struct {
	Store store.Store
}

func (s *CheckService) StartCheck() chan uptimemonitor.Monitor {
	ch := make(chan uptimemonitor.Monitor, 10)

	go func() {
		for m := range ch {
			fmt.Printf("#%d ", m.ID)
			s.handleCheck(m)
		}
	}()

	return ch
}

func (s *CheckService) RunCheck(ctx context.Context, ch chan uptimemonitor.Monitor) error {
	monitors, err := s.Store.ListMonitors(ctx)
	if err != nil {
		return err
	}

	for _, m := range monitors {
		fmt.Printf("x")
		ch <- m
	}

	fmt.Printf("\n")

	return nil
}

func (s *CheckService) handleCheck(m uptimemonitor.Monitor) {
	fmt.Print(".")

	c, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	start := time.Now()

	// todo: add timeout
	res, err := http.Get(m.Url)
	elapsed := time.Since(start)

	if err != nil {
		s.Store.CreateCheck(c, uptimemonitor.Check{
			MonitorID:      m.ID,
			Monitor:        m,
			StatusCode:     http.StatusInternalServerError,
			ResponseTimeMs: elapsed.Milliseconds(),
		})
		return
	}

	_, err = s.Store.CreateCheck(c, uptimemonitor.Check{
		MonitorID:      m.ID,
		Monitor:        m,
		StatusCode:     res.StatusCode,
		ResponseTimeMs: elapsed.Milliseconds(),
	})
	if err != nil {
		log.Printf("err: #%v", err)
		return
	}
}
