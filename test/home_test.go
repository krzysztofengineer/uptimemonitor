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
}
