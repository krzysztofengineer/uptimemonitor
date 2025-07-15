package test

import (
	"net/http"
	"testing"
	"time"
	"uptimemonitor"
	"uptimemonitor/handler"
)

func TestCheck_ListChecks(t *testing.T) {
	t.Run("setup is required to load checks", func(t *testing.T) {
		tc := NewTestCase(t)
		defer tc.Close()

		tc.Get("/monitors/1/checks").
			AssertRedirect(http.StatusSeeOther, "/setup")
	})

	t.Run("guests cannot load checks", func(t *testing.T) {
		tc := NewTestCase(t)
		defer tc.Close()

		tc.CreateTestUser("test@example.com", "password")

		tc.Get("/monitors/1/checks").
			AssertRedirect(http.StatusSeeOther, "/login")
	})

	t.Run("monitor has to exist", func(t *testing.T) {
		tc := NewTestCase(t)
		defer tc.Close()

		tc.LogIn().
			Get("/monitors/1/checks").
			AssertStatusCode(http.StatusNotFound)
	})

	t.Run("latest checks are returned", func(t *testing.T) {
		tc := NewTestCase(t)
		defer tc.Close()

		monitor, _ := tc.Store.CreateMonitor(t.Context(), uptimemonitor.Monitor{
			Url: "https://example.com",
		})

		tc.Store.CreateCheck(t.Context(), uptimemonitor.Check{
			MonitorID: monitor.ID,
			Monitor:   monitor,
		})

		tc.LogIn().
			Get("/monitors/1/checks").
			AssertStatusCode(http.StatusOK).
			AssertElementVisible(`div[id="monitors-1-checks-1"]`)
	})
}

func TestCheck_PeriodicChecks(t *testing.T) {
	t.Run("it does not create checks if no monitors are defined", func(t *testing.T) {
		tc := NewTestCase(t)
		defer tc.Close()

		handler := handler.CheckHandler{
			Store: tc.Store,
		}

		handler.RunCheck(t.Context())
		time.Sleep(time.Millisecond * 100)
		tc.AssertDatabaseCount("checks", 0)

		handler.RunCheck(t.Context())
		time.Sleep(time.Millisecond * 100)
		tc.AssertDatabaseCount("checks", 0)
	})

	t.Run("it creates checks every minute", func(t *testing.T) {
		tc := NewTestCase(t)
		defer tc.Close()

		tc.Store.CreateMonitor(t.Context(), uptimemonitor.Monitor{
			Url: "https://example.com",
		})

		handler := handler.CheckHandler{
			Store: tc.Store,
		}

		handler.RunCheck(t.Context())
		time.Sleep(time.Millisecond * 100)
		tc.AssertDatabaseCount("checks", 1)

		handler.RunCheck(t.Context())
		time.Sleep(time.Millisecond * 100)
		tc.AssertDatabaseCount("checks", 2)
	})
}
