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
}
