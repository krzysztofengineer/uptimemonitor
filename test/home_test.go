package test

import (
	"net/http"
	"testing"
)

func TestHome(t *testing.T) {
	t.Run("setup is required to access home page", func(t *testing.T) {
		tc := NewTestCase(t)
		defer tc.Close()

		tc.Get("/").AssertRedirect(http.StatusSeeOther, "/setup")
	})

	t.Run("guests cannot access home page", func(t *testing.T) {
		tc := NewTestCase(t)
		defer tc.Close()

		tc.CreateTestUser("test@example.com", "password")

		tc.Get("/").AssertRedirect(http.StatusSeeOther, "/login")
	})

	t.Run("monitors are displayed on home page", func(t *testing.T) {
		tc := NewTestCase(t)
		defer tc.Close()

		tc.LogIn().
			Get("/").
			AssertStatusCode(http.StatusOK).
			AssertElementVisible(`[hx-get="/monitors"]`)
	})

	t.Run("incidents are displayed on home page", func(t *testing.T) {
		tc := NewTestCase(t)
		defer tc.Close()

		tc.LogIn().
			Get("/").
			AssertStatusCode(http.StatusOK).
			AssertElementVisible(`[hx-get="/incidents"]`)
	})
}
