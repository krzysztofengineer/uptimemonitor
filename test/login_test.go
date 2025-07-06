package test

import (
	"net/http"
	"testing"
)

func TestLogin(t *testing.T) {
	t.Run("setup is required before user can log in", func(t *testing.T) {
		tc := NewTestCase(t)
		defer tc.Close()

		tc.Get("/login").AssertRedirect(http.StatusSeeOther, "/setup")
		tc.Post("/login", nil).AssertStatusCode(http.StatusForbidden)
	})
}
