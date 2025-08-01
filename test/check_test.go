package test

import (
	"net/http"
	"testing"
	"time"
	"uptimemonitor"
	"uptimemonitor/service"
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

		service := service.New(tc.Store)
		ch := service.StartCheck()

		service.RunCheck(t.Context(), ch)
		tc.AssertDatabaseCount("checks", 0)

		service.RunCheck(t.Context(), ch)
		tc.AssertDatabaseCount("checks", 0)
	})

	t.Run("it creates checks every minute", func(t *testing.T) {
		tc := NewTestCase(t)
		defer tc.Close()

		tc.Store.CreateMonitor(t.Context(), uptimemonitor.Monitor{
			Url: "https://example.com",
		})

		service := service.CheckService{
			Store: tc.Store,
		}
		ch := service.StartCheck()

		service.RunCheck(t.Context(), ch)
		service.RunCheck(t.Context(), ch)

		time.Sleep(3 * time.Second)
		// tc.AssertDatabaseCount("checks", 2) // todo: fix
	})

	t.Run("checks can use different http methods", func(t *testing.T) {
		tc := NewTestCase(t)
		defer tc.Close()

		tc.Store.CreateMonitor(t.Context(), uptimemonitor.Monitor{Url: tc.Server.URL + "/test/post", HttpMethod: http.MethodPost})
		tc.Store.CreateMonitor(t.Context(), uptimemonitor.Monitor{Url: tc.Server.URL + "/test/patch", HttpMethod: http.MethodPatch})
		tc.Store.CreateMonitor(t.Context(), uptimemonitor.Monitor{Url: tc.Server.URL + "/test/put", HttpMethod: http.MethodPut})
		tc.Store.CreateMonitor(t.Context(), uptimemonitor.Monitor{Url: tc.Server.URL + "/test/delete", HttpMethod: http.MethodDelete})

		service := service.CheckService{Store: tc.Store}
		ch := service.StartCheck()
		service.RunCheck(t.Context(), ch)
		time.Sleep(1 * time.Second)

		tc.AssertDatabaseCount("checks", 4)
		tc.AssertDatabaseCount("incidents", 0)

		first, err := tc.Store.GetCheckByID(t.Context(), 1)
		tc.AssertNoError(err)

		second, err := tc.Store.GetCheckByID(t.Context(), 2)
		tc.AssertNoError(err)

		third, err := tc.Store.GetCheckByID(t.Context(), 3)
		tc.AssertNoError(err)

		fifth, err := tc.Store.GetCheckByID(t.Context(), 4)
		tc.AssertNoError(err)

		tc.AssertEqual(http.StatusOK, first.StatusCode)
		tc.AssertEqual(http.StatusOK, second.StatusCode)
		tc.AssertEqual(http.StatusOK, third.StatusCode)
		tc.AssertEqual(http.StatusOK, fifth.StatusCode)
	})

	t.Run("checks can send custom body", func(t *testing.T) {
		tc := NewTestCase(t)
		defer tc.Close()

		tc.Store.CreateMonitor(t.Context(), uptimemonitor.Monitor{Url: tc.Server.URL + "/test/body", HttpMethod: http.MethodPost, HttpBody: `{"test":123}`})

		service := service.CheckService{Store: tc.Store}
		ch := service.StartCheck()
		service.RunCheck(t.Context(), ch)
		time.Sleep(1 * time.Second)

		check, _ := tc.Store.GetCheckByID(t.Context(), 1)
		tc.AssertEqual(http.StatusOK, check.StatusCode)
	})

	t.Run("checks can send custom headers", func(t *testing.T) {
		tc := NewTestCase(t)
		defer tc.Close()

		tc.Store.CreateMonitor(t.Context(), uptimemonitor.Monitor{Url: tc.Server.URL + "/test/headers", HttpMethod: http.MethodPost, HttpHeaders: `{"test":"abc"}`})

		service := service.CheckService{Store: tc.Store}
		ch := service.StartCheck()
		service.RunCheck(t.Context(), ch)
		time.Sleep(1 * time.Second)

		check, _ := tc.Store.GetCheckByID(t.Context(), 1)
		tc.AssertEqual(http.StatusOK, check.StatusCode)
	})
}

func TestCheck_Cleanup(t *testing.T) {
	t.Run("old cleanups are removed", func(t *testing.T) {
		tc := NewTestCase(t)
		defer tc.Close()

		tc.Store.CreateMonitor(t.Context(), uptimemonitor.Monitor{Url: tc.Server.URL + "/test/200", HttpMethod: http.MethodGet})

		tc.Store.CreateCheck(t.Context(), uptimemonitor.Check{
			MonitorID: 1,
			CreatedAt: time.Now().Add(-time.Hour).Add(-15 * time.Minute),
		})

		tc.Store.CreateCheck(t.Context(), uptimemonitor.Check{
			MonitorID: 1,
		})

		tc.Store.CreateIncident(t.Context(), uptimemonitor.Incident{
			MonitorID: 1,
			CreatedAt: time.Now().Add(-time.Hour * 24 * 8),
		})

		tc.Store.CreateIncident(t.Context(), uptimemonitor.Incident{
			MonitorID: 1,
		})

		tc.AssertDatabaseCount("checks", 2)
		tc.AssertDatabaseCount("incidents", 2)

		service := service.CheckService{Store: tc.Store}
		ch := service.StartCheck()
		service.RunCheck(t.Context(), ch)
		time.Sleep(1 * time.Second)

		tc.AssertDatabaseCount("checks", 2) // one new, old one deleted
		tc.AssertDatabaseCount("incidents", 1)

		service.RunCheck(t.Context(), ch)
		time.Sleep(1 * time.Second)

		tc.AssertDatabaseCount("checks", 3)
	})
}
