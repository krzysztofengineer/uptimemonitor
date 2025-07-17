package service

import (
	"context"
	"fmt"
	"log"
	"time"
	"uptimemonitor"
	"uptimemonitor/store"

	"math/rand/v2"
)

type CheckService struct {
	Store store.Store
}

func NewCheckService(store store.Store) *CheckService {
	return &CheckService{
		Store: store,
	}
}

func (s *CheckService) StartCheck() chan uptimemonitor.Monitor {
	ch := make(chan uptimemonitor.Monitor, 10)

	go func() {
		for m := range ch {
			s.handleCheck(m)
		}
	}()

	return ch
}

func (s *CheckService) RunChecks(ctx context.Context, ch chan uptimemonitor.Monitor) error {
	monitors, err := s.Store.ListMonitors(ctx)
	if err != nil {
		return err
	}

	for _, m := range monitors {
		ch <- m
	}

	return nil
}

func (s *CheckService) handleCheck(m uptimemonitor.Monitor) {
	fmt.Print(".")

	c, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	time.Sleep(time.Duration(rand.IntN(500)) * time.Millisecond)

	_, err := s.Store.CreateCheck(c, uptimemonitor.Check{
		MonitorID: m.ID,
		Monitor:   m,
	})
	if err != nil {
		log.Printf("err: #%v", err)
		return
	}
}
