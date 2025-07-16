package service

import (
	"context"
	"log"
	"time"
	"uptimemonitor"
	"uptimemonitor/store"
)

type CheckService struct {
	Store store.Store

	channel chan uptimemonitor.Monitor
	Done    chan bool
}

func NewCheckService(store store.Store) *CheckService {
	return &CheckService{
		Store:   store,
		channel: make(chan uptimemonitor.Monitor),
		Done:    make(chan bool),
	}
}

func (s *CheckService) Start() {
	for {
		select {
		case <-s.Done:
			return
		case m := <-s.channel:
			log.Printf("got m: %v", m)
			s.handleCheck(m)
		}
	}
}

func (s *CheckService) RunChecks(ctx context.Context) error {
	monitors, err := s.Store.ListMonitors(ctx)
	if err != nil {
		return err
	}

	log.Printf("running check: %d", len(monitors))

	for _, m := range monitors {
		s.channel <- m
	}

	return nil
}

func (s *CheckService) handleCheck(m uptimemonitor.Monitor) {
	c, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	log.Printf("CHECK #%d", m.ID)

	check, err := s.Store.CreateCheck(c, uptimemonitor.Check{
		MonitorID: m.ID,
		Monitor:   m,
	})
	if err != nil {
		log.Printf("err: #%v", err)
		return
	}

	log.Printf("CHECK FINISHED WITH ID: #%d", check.ID)
}
