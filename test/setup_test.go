package test

import (
	"net/http"
	"testing"
	"time"
	"uptimemonitor"
)

func TestSetup(t *testing.T) {
	t.Run("redirects to setup when no users are found", func(t *testing.T) {
		tc := NewTestCase(t)
		defer tc.Close()

		tc.Get("/").
			AssertRedirect(http.StatusSeeOther, "/setup")
	})

	t.Run("redirects to home page when users are found", func(t *testing.T) {
		tc := NewTestCase(t)
		defer tc.Close()

		tc.Store.CreateUser(t.Context(), uptimemonitor.User{
			Name:      "Test User",
			Email:     "test@example.com",
			CreatedAt: time.Now(),
		})

		tc.Get("/setup").
			AssertRedirect(http.StatusSeeOther, "/")
	})

	t.Run("shows a setup form when no users are found", func(t *testing.T) {
		tc := NewTestCase(t)
		defer tc.Close()

		tc.Get("/setup").
			AssertStatusCode(http.StatusOK).
			AssertElementVisible(`form[hx-post="/setup"]`).
			AssertElementVisible(`input[name="name"]`).
			AssertElementVisible(`input[name="email"]`).
			AssertElementVisible(`input[name="password"]`).
			AssertElementVisible(`button[type="submit"]`)
	})
}
