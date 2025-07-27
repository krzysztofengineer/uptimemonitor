package test

import (
	"fmt"
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

	t.Run("invalid domains create incidents", func(t *testing.T) {
		tc := NewTestCase(t)
		defer tc.Close()

		service := service.New(tc.Store)

		tc.Store.CreateMonitor(t.Context(), uptimemonitor.Monitor{
			Url: "http://invalid-url",
		})

		ch := service.StartCheck()
		service.RunCheck(t.Context(), ch)

		time.Sleep(time.Second)
		tc.AssertDatabaseCount("incidents", 1)
	})

	t.Run("incidents get resolved", func(t *testing.T) {
		tc := NewTestCase(t)
		defer tc.Close()

		service := service.New(tc.Store)

		tc.Store.CreateMonitor(t.Context(), uptimemonitor.Monitor{
			Url: tc.Server.URL + "/test/even",
		})

		ch := service.StartCheck()
		service.RunCheck(t.Context(), ch)

		time.Sleep(time.Second)
		tc.AssertDatabaseCount("incidents", 1)

		service.RunCheck(t.Context(), ch)
		time.Sleep(time.Second)
		tc.AssertDatabaseCount("incidents", 1)

		incident, _ := tc.Store.LastOpenIncident(t.Context(), 1)
		if incident.ID != 0 {
			t.Fatalf("expected not to found any incidents,found: %v", incident)
		}
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

func TestIncident_ListMonitorIncidents(t *testing.T) {
	t.Run("setup is required", func(t *testing.T) {
		tc := NewTestCase(t)
		defer tc.Close()

		tc.Get("/monitors/1/incidents").AssertRedirect(http.StatusSeeOther, "/setup")
	})

	t.Run("guests cannot list incidents", func(t *testing.T) {
		tc := NewTestCase(t)
		defer tc.Close()

		tc.CreateTestUser("test@example.com", "password")

		tc.Get("/monitors/1/incidents").AssertRedirect(http.StatusSeeOther, "/login")
	})

	t.Run("users can list incidents", func(t *testing.T) {
		tc := NewTestCase(t)
		defer tc.Close()

		tc.Store.CreateMonitor(t.Context(), uptimemonitor.Monitor{
			Url: "http://example.com",
		})

		tc.Store.CreateIncident(t.Context(), uptimemonitor.Incident{
			MonitorID:      1,
			StatusCode:     404,
			ResponseTimeMs: 100,
		})

		tc.LogIn().Get("/monitors/1/incidents").
			AssertStatusCode(http.StatusOK).
			AssertElementVisible(`[id="incidents-1"]`)
	})
}

func TestIncident_RemoveIncidents(t *testing.T) {
	t.Run("setup is required", func(t *testing.T) {
		tc := NewTestCase(t)
		defer tc.Close()

		tc.Delete("/incidents/1").AssertStatusCode(http.StatusForbidden)
	})

	t.Run("guests cannot remove incidents", func(t *testing.T) {
		tc := NewTestCase(t)
		defer tc.Close()

		tc.CreateTestUser("test@example.com", "password")
		tc.Delete("/incidents/1").AssertRedirect(http.StatusSeeOther, "/login")
	})

	t.Run("remove incident form is visible", func(t *testing.T) {
		tc := NewTestCase(t)
		defer tc.Close()

		m, _ := tc.Store.CreateMonitor(t.Context(), uptimemonitor.Monitor{
			Url: "http://example.com",
		})

		i, _ := tc.Store.CreateIncident(t.Context(), uptimemonitor.Incident{
			MonitorID:      1,
			StatusCode:     404,
			ResponseTimeMs: 100,
		})

		tc.LogIn().
			Get(fmt.Sprintf("/m/%s/i/%s", m.Uuid, i.Uuid)).
			AssertStatusCode(http.StatusOK).
			AssertElementVisible(`form[hx-delete="/incidents/1"]`)
	})

	t.Run("users can remove incidents", func(t *testing.T) {
		tc := NewTestCase(t)
		defer tc.Close()

		m, _ := tc.Store.CreateMonitor(t.Context(), uptimemonitor.Monitor{
			Url: "http://example.com",
		})

		tc.Store.CreateIncident(t.Context(), uptimemonitor.Incident{
			MonitorID:      1,
			StatusCode:     404,
			ResponseTimeMs: 100,
		})

		tc.LogIn().
			Delete("/incidents/1").
			AssertStatusCode(http.StatusOK).
			AssertHeader("HX-Redirect", fmt.Sprintf("/m/%s", m.Uuid))

		tc.AssertDatabaseCount("incidents", 0)
	})
}

func TestIncident_IncidentPage(t *testing.T) {
	t.Run("setup is required", func(t *testing.T) {
		tc := NewTestCase(t)
		defer tc.Close()

		tc.Get("/m/uuid/i/uuid").AssertRedirect(http.StatusSeeOther, "/setup")
	})

	t.Run("guests cannot view incidents", func(t *testing.T) {
		tc := NewTestCase(t)
		defer tc.Close()

		tc.CreateTestUser("test@example.com", "password")

		tc.Get("/m/uuid/i/uuid").AssertRedirect(http.StatusSeeOther, "/login")

	})

	t.Run("incident has to exist", func(t *testing.T) {
		tc := NewTestCase(t)
		defer tc.Close()

		m, _ := tc.Store.CreateMonitor(t.Context(), uptimemonitor.Monitor{
			Url: "http://example.com",
		})

		tc.LogIn().Get(fmt.Sprintf("/m/%s/i/uuid", m.Uuid)).
			AssertStatusCode(http.StatusNotFound)
	})

	t.Run("monitor has to exist", func(t *testing.T) {
		tc := NewTestCase(t)
		defer tc.Close()
		tc.Store.CreateMonitor(t.Context(), uptimemonitor.Monitor{
			Url: "http://example.com",
		})

		i, _ := tc.Store.CreateIncident(t.Context(), uptimemonitor.Incident{
			MonitorID:      1,
			StatusCode:     404,
			ResponseTimeMs: 100,
		})

		tc.LogIn().Get(fmt.Sprintf("/m/uuid/i/%s", i.Uuid)).
			AssertStatusCode(http.StatusNotFound)
	})

	t.Run("incident is visible", func(t *testing.T) {
		tc := NewTestCase(t)
		defer tc.Close()

		m, _ := tc.Store.CreateMonitor(t.Context(), uptimemonitor.Monitor{
			Url: "http://example.com",
		})

		i, _ := tc.Store.CreateIncident(t.Context(), uptimemonitor.Incident{
			MonitorID:      1,
			StatusCode:     404,
			ResponseTimeMs: 100,
		})

		tc.LogIn().
			Get(fmt.Sprintf("/m/%s/i/%s", m.Uuid, i.Uuid)).
			AssertStatusCode(http.StatusOK).
			AssertSeeText("404 Not Found")
	})

}
