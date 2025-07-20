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
			AssertSeeText("The email is required").
			AssertSeeText("The password is required")

		tc.Post("/login", url.Values{
			"email": []string{"invalid"},
		}).
			AssertStatusCode(http.StatusBadRequest).
			AssertSeeText("The email format is invalid")
	})

	t.Run("user has to exist", func(t *testing.T) {
		tc := NewTestCase(t)
		defer tc.Close()

		tc.CreateTestUser("test@example.com", "password")

		tc.Post("/login", url.Values{
			"email":    []string{"other@example.com"},
			"password": []string{"password"},
		}).
			AssertStatusCode(http.StatusBadRequest).
			AssertSeeText("The credentials do not match our records")
	})

	t.Run("password has to be valid", func(t *testing.T) {
		tc := NewTestCase(t)
		defer tc.Close()

		tc.CreateTestUser("test@example.com", "password")

		tc.Post("/login", url.Values{
			"email":    []string{"test@example.com"},
			"password": []string{"invalid"},
		}).
			AssertStatusCode(http.StatusBadRequest).
			AssertSeeText("The credentials do not match our records")
	})

	t.Run("user can log in", func(t *testing.T) {
		tc := NewTestCase(t)
		defer tc.Close()

		tc.CreateTestUser("test@example.com", "password")

		tc.Get("/").AssertRedirect(http.StatusSeeOther, "/login")

		ar := tc.Post("/login", url.Values{
			"email":    []string{"test@example.com"},
			"password": []string{"password"},
		}).
			AssertStatusCode(http.StatusOK).
			AssertCookieSet("session").
			AssertHeader("HX-Redirect", "/")

		tc.AssertDatabaseCount("sessions", 1)

		cookies := ar.Response.Cookies()
		var cookie *http.Cookie

		for _, c := range cookies {
			if c.Name == "session" {
				cookie = c
			}
		}

		tc.WithCookie(cookie).
			Get("/new").
			AssertNoRedirect().
			AssertStatusCode(http.StatusOK)
	})

	t.Run("logged in users are redirected to a new page", func(t *testing.T) {
		tc := NewTestCase(t)
		defer tc.Close()

		tc.LogIn().
			Get("/login").
			AssertRedirect(http.StatusSeeOther, "/new")
	})
}
