package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"text/template"
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

	go func() {
		s.Cleanup()
	}()

	return nil
}

func (s *CheckService) handleCheck(m uptimemonitor.Monitor) {
	c, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	start := time.Now()
	var customBody io.Reader

	if m.HttpBody != "" {
		customBody = strings.NewReader(m.HttpBody)
	}

	req, err := http.NewRequest(
		m.HttpMethod,
		m.Url,
		customBody,
	)

	if err != nil {
		return
	}

	if m.HttpHeaders != "" {
		customHeaders := map[string]string{}
		err = json.Unmarshal([]byte(m.HttpHeaders), &customHeaders)
		if err == nil {
			for k, v := range customHeaders {
				req.Header.Add(k, v)
			}
		}
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
	reqURL := m.Url
	reqMethod := m.HttpMethod
	reqHeaders := m.HttpHeaders
	reqBody := m.HttpBody

	if exists, latest := s.incidentAlreadyExists(context.Background(), m, statusCode); exists {
		return s.Store.UpdateIncidentBodyAndHeaders(context.Background(), latest, body, headers, reqMethod, reqURL, reqHeaders, reqBody)
	}

	s.Store.ResolveMonitorIncidents(context.Background(), m)

	incident := uptimemonitor.Incident{
		MonitorID:      m.ID,
		StatusCode:     check.StatusCode,
		ResponseTimeMs: responseTimeMs,
		Body:           body,
		Headers:        headers,
		ReqUrl:         reqURL,
		ReqHeaders:     reqHeaders,
		ReqBody:        reqBody,
		ReqMethod:      reqMethod,
	}

	saved, err := s.Store.CreateIncident(context.Background(), incident)
	if err != nil {
		return fmt.Errorf("failed to create incident for monitor %d: %w", m.ID, err)
	}

	if m.WebhookUrl != "" {
		s.callWebhook(m, saved)
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

func (s *CheckService) callWebhook(m uptimemonitor.Monitor, i uptimemonitor.Incident) {
	var customBody io.Reader

	// todo: parse url
	if m.WebhookBody != "" {
		t, err := template.New("webhook").Parse(m.WebhookBody)
		if err != nil {
			customBody = strings.NewReader(m.WebhookBody)
		} else {
			var buf bytes.Buffer
			err = t.Execute(&buf, struct {
				Url        string
				StatusCode int
			}{
				Url:        m.Url,
				StatusCode: i.StatusCode,
			})
			if err != nil {
				customBody = strings.NewReader(m.WebhookBody)
			} else {
				customBody = strings.NewReader(buf.String())
			}
		}

	}

	req, err := http.NewRequest(
		m.WebhookMethod,
		m.WebhookUrl,
		customBody,
	)

	if err != nil {
		return
	}

	if m.WebhookHeaders != "" {
		customHeaders := map[string]string{}
		err = json.Unmarshal([]byte(m.WebhookHeaders), &customHeaders)
		if err == nil {
			for k, v := range customHeaders {
				req.Header.Add(k, v)
			}
		}
	}

	http.DefaultClient.Do(req)
}

func (s *CheckService) Cleanup() {
	s.Store.DeleteOldChecks(context.Background())
	s.Store.DeleteOldIncidents(context.Background())
}
