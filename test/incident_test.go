package test

import (
	"net/http"
	"testing"
	"time"
	"uptimemonitor"
	"uptimemonitor/service"
)

func TestIncident(t *testing.T) {
	t.Run("no incident is created when check succeeds", func(t *testing.T) {
		tc := NewTestCase(t)
		defer tc.Close()

		service := service.New(tc.Store)

		tc.Store.CreateMonitor(t.Context(), uptimemonitor.Monitor{
			Url: tc.Server.URL + "/test/200",
		})

		ch := service.StartCheck()
		service.RunCheck(t.Context(), ch)

		time.Sleep(time.Second)

		tc.AssertDatabaseCount("incidents", 0)
		tc.AssertDatabaseCount("checks", 1)
	})

	t.Run("incident is created when check fails", func(t *testing.T) {
		tc := NewTestCase(t)
		defer tc.Close()

		service := service.New(tc.Store)

		tc.Store.CreateMonitor(t.Context(), uptimemonitor.Monitor{
			Url: tc.Server.URL + "/test/404",
		})

		ch := service.StartCheck()
		service.RunCheck(t.Context(), ch)

		time.Sleep(time.Second)

		tc.AssertDatabaseCount("incidents", 1)
		tc.AssertDatabaseCount("checks", 1)
	})

	t.Run("new incident is not created for the same monitor if it already exists", func(t *testing.T) {
		tc := NewTestCase(t)
		defer tc.Close()

		service := service.New(tc.Store)

		tc.Store.CreateMonitor(t.Context(), uptimemonitor.Monitor{
			Url: tc.Server.URL + "/test/500",
		})

		ch := service.StartCheck()
		service.RunCheck(t.Context(), ch)

		time.Sleep(time.Second)

		tc.AssertDatabaseCount("incidents", 1)
		tc.AssertDatabaseCount("checks", 1)

		ch = service.StartCheck()
		service.RunCheck(t.Context(), ch)

		time.Sleep(time.Second)

		tc.AssertDatabaseCount("incidents", 1)
	})

	t.Run("incident is created when check fails with different status code", func(t *testing.T) {
		tc := NewTestCase(t)
		defer tc.Close()

		service := service.New(tc.Store)

		tc.Store.CreateMonitor(t.Context(), uptimemonitor.Monitor{
			Url: tc.Server.URL + "/test/500",
		})

		tc.Store.CreateIncident(t.Context(), uptimemonitor.Incident{
			MonitorID:      1,
			StatusCode:     404,
			ResponseTimeMs: 100,
			Body:           "not found",
			Headers:        "Content-Type: text/plain",
		})

		ch := service.StartCheck()
		service.RunCheck(t.Context(), ch)

		time.Sleep(time.Second)

		tc.AssertDatabaseCount("incidents", 2)
	})
}

func TestIncident_ListIncidents(t *testing.T) {
	t.Run("setup is required", func(t *testing.T) {
		tc := NewTestCase(t)
		defer tc.Close()

		tc.Get("/incidents").AssertRedirect(http.StatusSeeOther, "/setup")
	})

	t.Run("guests cannot list incidents", func(t *testing.T) {
		tc := NewTestCase(t)
		defer tc.Close()

		tc.CreateTestUser("test@example.com", "password")

		tc.Get("/incidents").AssertRedirect(http.StatusSeeOther, "/login")
	})

	t.Run("users can list incidents", func(t *testing.T) {
		tc := NewTestCase(t)
		defer tc.Close()

		tc.Store.CreateMonitor(t.Context(), uptimemonitor.Monitor{
			Url: "http://example.com",
		})
		tc.Store.CreateMonitor(t.Context(), uptimemonitor.Monitor{
			Url: "http://example.com",
		})

		tc.Store.CreateIncident(t.Context(), uptimemonitor.Incident{
			MonitorID:      1,
			StatusCode:     404,
			ResponseTimeMs: 100,
		})

		tc.Store.CreateIncident(t.Context(), uptimemonitor.Incident{
			MonitorID:      2,
			StatusCode:     500,
			ResponseTimeMs: 100,
		})

		tc.LogIn().Get("/incidents").
			AssertStatusCode(http.StatusOK).
			AssertElementVisible(`[id="incidents-1"]`).
			AssertElementVisible(`[id="incidents-2"]`)
	})
}
