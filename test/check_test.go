package test

import (
	"net/http"
	"sync"
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

		var wg sync.WaitGroup
		wg.Add(2)

		service := service.CheckService{
			Store: tc.Store,
		}
		ch := service.StartCheck()

		go func() {
			defer wg.Done()
			service.RunCheck(t.Context(), ch)
		}()
		go func() {
			defer wg.Done()
			service.RunCheck(t.Context(), ch)
		}()

		wg.Wait()
		close(ch)
		time.Sleep(1 * time.Second)
		tc.AssertDatabaseCount("checks", 2)
	})
}
