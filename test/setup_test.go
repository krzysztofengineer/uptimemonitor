package test

import (
	"net/http"
	"net/url"
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

	t.Run("validates a form", func(t *testing.T) {
		tc := NewTestCase(t)
		defer tc.Close()

		tc.Post("/setup", url.Values{}).
			AssertStatusCode(http.StatusBadRequest).
			AssertElementVisible(`form[hx-swap="outerHTML"]`).
			AssertSeeText("The name field is required").
			AssertSeeText("The email field is required").
			AssertSeeText("The password field is required")

		res := tc.Post("/setup", url.Values{
			"email": []string{"invalid"},
		})

		res.AssertStatusCode(http.StatusBadRequest).
			AssertElementVisible(`form[hx-swap="outerHTML"]`).
			AssertSeeText("The email format is invalid")
	})

	t.Run("setup", func(t *testing.T) {
		tc := NewTestCase(t)
		defer tc.Close()

		tc.AssertDatabaseCount("users", 0)

		res := tc.Post("/setup", url.Values{
			"name":     []string{"Test"},
			"email":    []string{"test@example.com"},
			"password": []string{"password"},
		})

		res.AssertHeader("HX-Redirect", "/")
		tc.AssertDatabaseCount("users", 1)

		tc.Get("/setup").AssertRedirect(http.StatusSeeOther, "/")
	})
}
