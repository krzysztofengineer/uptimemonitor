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
}

func NewCheckService(store store.Store) *CheckService {
	return &CheckService{
		Store: store,
	}
}

func (s *Service) RunChecks(ctx context.Context) error {
	monitors, err := s.Store.ListMonitors(ctx)
	if err != nil {
		return err
	}

	log.Printf("running check: %d", len(monitors))

	for _, m := range monitors {
		go func(mon uptimemonitor.Monitor) {
			c, cancel := context.WithTimeout(context.Background(), time.Minute)
			defer cancel()

			log.Printf("CHECK #%d", m.ID)

			check, err := s.Store.CreateCheck(c, uptimemonitor.Check{
				MonitorID: mon.ID,
				Monitor:   mon,
			})
			if err != nil {
				log.Printf("err: #%v", err)
				return
			}

			log.Printf("CHECK FINISHED WITH ID: #%d", check.ID)
		}(m)
	}

	return nil
}
