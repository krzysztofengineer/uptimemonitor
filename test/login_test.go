package test

import (
	"net/http"
	"net/url"
	"testing"
)

func TestLogin(t *testing.T) {
	t.Run("setup is required before user can log in", func(t *testing.T) {
		tc := NewTestCase(t)
		defer tc.Close()

		tc.Get("/login").AssertRedirect(http.StatusSeeOther, "/setup")
		tc.Post("/login", nil).AssertStatusCode(http.StatusForbidden)
	})

	t.Run("shows a login form", func(t *testing.T) {
		tc := NewTestCase(t)
		defer tc.Close()

		tc.CreateTestUser("test@example.com", "password")
		tc.Get("/login").
			AssertStatusCode(http.StatusOK).
			AssertElementVisible(`form[hx-post="/login"]`)
	})

	t.Run("validates form", func(t *testing.T) {
		tc := NewTestCase(t)
		defer tc.Close()

		tc.CreateTestUser("test@example.com", "password")

		tc.Post("/login", nil).
			AssertStatusCode(http.StatusBadRequest).
			AssertSeeText("The email field is required").
			AssertSeeText("The password field is required")

		tc.Post("/login", url.Values{
			"email": []string{"invalid"},
		}).
			AssertStatusCode(http.StatusBadRequest).
			AssertSeeText("The email format is invalid")
	})
}
