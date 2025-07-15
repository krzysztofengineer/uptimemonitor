package test

import (
	"net/http"
	"testing"
	"uptimemonitor"
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
