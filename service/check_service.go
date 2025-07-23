package service

import (
	"context"
	"fmt"
	"io"
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
		ch <- m
	}

	return nil
}

func (s *CheckService) handleCheck(m uptimemonitor.Monitor) {
	c, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	start := time.Now()

	req, err := http.NewRequest(
		m.HttpMethod,
		m.Url,
		nil,
	)

	if err != nil {
		return
	}

	// todo: add timeout
	res, err := http.DefaultClient.Do(req)
	elapsed := time.Since(start)

	if err != nil {
		check, err := s.Store.CreateCheck(c, uptimemonitor.Check{
			MonitorID:      m.ID,
			Monitor:        m,
			StatusCode:     http.StatusNotFound,
			ResponseTimeMs: elapsed.Milliseconds(),
		})

		if err != nil || check.ID == 0 {
			return
		}

		if err := s.createIncident(m, check, elapsed.Milliseconds(), http.StatusNotFound, "", ""); err != nil {
			return
		}

		return
	}

	statusCode := res.StatusCode
	if res.Request.Response != nil {
		statusCode = res.Request.Response.StatusCode
	}

	check, err := s.Store.CreateCheck(c, uptimemonitor.Check{
		MonitorID:      m.ID,
		Monitor:        m,
		StatusCode:     statusCode,
		ResponseTimeMs: elapsed.Milliseconds(),
	})

	if err != nil {
		return
	}

	if statusCode < 300 {
		open, err := s.Store.ListMonitorOpenIncidents(c, m.ID)
		if err != nil {
			return
		}

		for _, i := range open {
			s.Store.ResolveIncident(c, i)
		}

		return
	}

	var body string
	var headers string

	if res.Request.Response != nil {
		var resHeaders string
		for k, v := range res.Request.Response.Header {
			resHeaders += fmt.Sprintf("%s: %s\n", k, v)
		}

		resBody, _ := io.ReadAll(res.Request.Response.Body)
		defer res.Request.Response.Body.Close()
		body = string(resBody)
		headers = resHeaders
	} else {
		var resHeaders string
		for k, v := range res.Header {
			resHeaders += fmt.Sprintf("%s: %s\n", k, v)
		}

		resBody, _ := io.ReadAll(res.Body)
		defer res.Body.Close()
		body = string(resBody)
		headers = resHeaders
	}

	s.createIncident(m, check, elapsed.Milliseconds(), statusCode, string(body), headers)
}

func (s *CheckService) createIncident(m uptimemonitor.Monitor, check uptimemonitor.Check, responseTimeMs int64, statusCode int, body string, headers string) error {
	if exists, latest := s.incidentAlreadyExists(context.Background(), m, statusCode); exists {
		return s.Store.UpdateIncidentBodyAndHeaders(context.Background(), latest, body, headers)
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

func (s *CheckService) incidentAlreadyExists(ctx context.Context, m uptimemonitor.Monitor, statusCode int) (bool, uptimemonitor.Incident) {
	latest, err := s.Store.LastIncidentByStatusCode(ctx, m.ID, uptimemonitor.IncidentStatusOpen, statusCode)
	if err != nil {
		return false, uptimemonitor.Incident{}
	}

	return latest.StatusCode == statusCode, latest
}
