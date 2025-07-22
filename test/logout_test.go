package test

import (
	"net/http"
	"testing"
)

func TestLogout(t *testing.T) {
	t.Run("user can log out", func(t *testing.T) {
		tc := NewTestCase(t)
		defer tc.Close()

		tc.LogIn().
			Get("/logout").
			AssertStatusCode(http.StatusOK).
			AssertRedirect(http.StatusSeeOther, "/login")

		tc.AssertDatabaseCount("sessions", 0)

		tc.Get("/").AssertRedirect(http.StatusSeeOther, "/login")
	})
}
