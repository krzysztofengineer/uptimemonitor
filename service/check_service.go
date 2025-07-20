package service

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
	"uptimemonitor"
	"uptimemonitor/store"
)

type CheckService struct {
	Store *store.Store
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
		check, err := s.Store.CreateCheck(c, uptimemonitor.Check{
			MonitorID:      m.ID,
			Monitor:        m,
			StatusCode:     http.StatusInternalServerError,
			ResponseTimeMs: elapsed.Milliseconds(),
		})
		if err != nil || check.ID == 0 {
			log.Printf("failed to create check: %v", err)
			return
		}

		if err := s.createIncident(m, check, elapsed.Milliseconds(), http.StatusInternalServerError, "", ""); err != nil {
			log.Printf("failed to create incident: %v", err)
		}

		return
	}

	check, err := s.Store.CreateCheck(c, uptimemonitor.Check{
		MonitorID:      m.ID,
		Monitor:        m,
		StatusCode:     res.StatusCode,
		ResponseTimeMs: elapsed.Milliseconds(),
	})

	if err != nil {
		log.Printf("err: #%v", err)
		return
	}

	if res.StatusCode >= 400 {
		body, err := io.ReadAll(res.Body)
		defer res.Body.Close()

		if err != nil {
			s.createIncident(m, check, elapsed.Milliseconds(), res.StatusCode, "", fmt.Sprintf("%v", res.Header))
			return
		}

		s.createIncident(m, check, elapsed.Milliseconds(), res.StatusCode, string(body), fmt.Sprintf("%v", res.Header))
	}
}

func (s *CheckService) createIncident(m uptimemonitor.Monitor, check uptimemonitor.Check, responseTimeMs int64, statusCode int, body string, headers string) error {
	if s.incidentAlreadyExists(context.Background(), m, statusCode) {
		return nil
	}

	incident := uptimemonitor.Incident{
		MonitorID:      m.ID,
		StatusCode:     check.StatusCode,
		ResponseTimeMs: responseTimeMs,
		Body:           body,
		Headers:        headers,
	}

	if _, err := s.Store.CreateIncident(context.Background(), incident); err != nil {
		return fmt.Errorf("failed to create incident for monitor %d: %w", m.ID, err)
	}

	return nil
}

func (s *CheckService) incidentAlreadyExists(ctx context.Context, m uptimemonitor.Monitor, statusCode int) bool {
	latest, err := s.Store.LastIncident(ctx, m.ID, uptimemonitor.IncidentStatusOpen, statusCode)
	if err != nil {
		return false
	}

	return latest.StatusCode == statusCode
}
