package test

import (
	"net/http"
	"net/url"
	"testing"
	"time"
	"uptimemonitor"
)

func TestMonitor_ListMonitors(t *testing.T) {
	t.Run("setup is required", func(t *testing.T) {
		tc := NewTestCase(t)
		defer tc.Close()

		tc.Get("/monitors").
			AssertRedirect(http.StatusSeeOther, "/setup")
	})

	t.Run("logged user is required", func(t *testing.T) {
		tc := NewTestCase(t)
		defer tc.Close()

		tc.CreateTestUser("test@example.com", "password")

		tc.Get("/monitors").
			AssertRedirect(http.StatusSeeOther, "/login")
	})

	t.Run("monitors table is visible on home page", func(t *testing.T) {
		tc := NewTestCase(t)
		defer tc.Close()

		tc.LogIn().
			Get("/").
			AssertStatusCode(http.StatusOK).
			AssertElementVisible(`div[hx-get="/monitors"]`)
	})

	t.Run("empty monitors list", func(t *testing.T) {
		tc := NewTestCase(t)
		defer tc.Close()

		tc.LogIn().
			Get("/monitors").
			AssertNoRedirect().
			AssertStatusCode(http.StatusOK)
	})

	t.Run("list monitors", func(t *testing.T) {
		tc := NewTestCase(t)
		defer tc.Close()

		tc.Store.CreateMonitor(t.Context(), uptimemonitor.Monitor{Url: "https://example.com", CreatedAt: time.Now()})
		tc.Store.CreateMonitor(t.Context(), uptimemonitor.Monitor{Url: "https://example.com/123", CreatedAt: time.Now()})

		tc.LogIn().
			Get("/monitors").
			AssertSeeText("example.com").
			AssertSeeText("example.com/123")
	})
}

func TestMonitor_CreateMonitor(t *testing.T) {
	t.Run("setup is required", func(t *testing.T) {
		tc := NewTestCase(t)
		defer tc.Close()

		tc.Post("/monitors", url.Values{}).
			AssertStatusCode(http.StatusForbidden)
	})

	t.Run("user has to be logged in", func(t *testing.T) {
		tc := NewTestCase(t)
		defer tc.Close()

		tc.CreateTestUser("test@example.com", "password")

		tc.Post("/monitors", url.Values{}).
			AssertRedirect(http.StatusSeeOther, "/login")
	})

	t.Run("url is required", func(t *testing.T) {
		tc := NewTestCase(t)
		defer tc.Close()

		tc.LogIn().
			Post("/monitors", url.Values{}).
			AssertNoRedirect().
			AssertStatusCode(http.StatusBadRequest).
			AssertSeeText("The url is required")
	})

	t.Run("the url has to be a valid url", func(t *testing.T) {
		tc := NewTestCase(t)
		defer tc.Close()

		tc.LogIn().
			Post("/monitors", url.Values{
				"url": []string{"invalid"},
			}).
			AssertStatusCode(http.StatusBadRequest).
			AssertSeeText("The url is invalid")
	})
}
