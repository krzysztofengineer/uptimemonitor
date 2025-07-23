package test

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"
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

		tc.Store.CreateMonitor(t.Context(), uptimemonitor.Monitor{Url: "https://example.com"})

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

		tc.Store.CreateMonitor(t.Context(), uptimemonitor.Monitor{Url: "https://example.com"})
		tc.Store.CreateMonitor(t.Context(), uptimemonitor.Monitor{Url: "https://example.com/123"})

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

	t.Run("monitor form is visible", func(t *testing.T) {
		tc := NewTestCase(t)
		defer tc.Close()

		tc.LogIn().
			Get("/new").
			AssertNoRedirect().
			AssertStatusCode(http.StatusOK).
			AssertElementVisible(`form[hx-post="/monitors"]`).
			AssertElementVisible(`select[name="http_method"]`).
			AssertElementVisible(`input[name="url"]`)
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

	t.Run("the url can be created", func(t *testing.T) {
		tc := NewTestCase(t)
		defer tc.Close()

		res := tc.LogIn().
			Post("/monitors", url.Values{
				"http_method":        []string{"GET"},
				"has_custom_headers": []string{"on"},
				"http_headers":       []string{`{"test":"abc"}`},
				"has_custom_body":    []string{"on"},
				"http_body":          []string{`{"test":"123"}`},
				"url":                []string{"https://example.com"},
			}).
			AssertStatusCode(http.StatusOK)

		m, _ := tc.Store.GetMonitorByID(t.Context(), 1)

		res.AssertHeader("HX-Redirect", fmt.Sprintf("/m/%s", m.Uuid))
		tc.AssertDatabaseCount("monitors", 1)
		tc.Get("/monitors").AssertSeeText("example.com")

		tc.AssertEqual(m.Url, "https://example.com")
		tc.AssertEqual(m.HttpMethod, "GET")
		tc.AssertEqual(m.HttpHeaders, `{"test":"abc"}`)
		tc.AssertEqual(m.HttpBody, `{"test":"123"}`)
	})

	t.Run("custom headers are validated when present", func(t *testing.T) {
		tc := NewTestCase(t)
		defer tc.Close()

		tc.LogIn().
			Post("/monitors", url.Values{
				"http_method": []string{"GET"},
				"url":         []string{"https://example.com"},
			}).
			AssertStatusCode(http.StatusOK)

		tc.Post("/monitors", url.Values{
			"http_method":        []string{"GET"},
			"url":                []string{"https://example.com"},
			"has_custom_headers": []string{"on"},
			"http_headers":       []string{`INVALID JSON`},
		}).
			AssertStatusCode(http.StatusBadRequest).
			AssertSeeText("The http headers should be a valid JSON")
	})
}

func TestMonitor_MonitorPage(t *testing.T) {
	t.Run("setup is required", func(t *testing.T) {
		tc := NewTestCase(t)
		defer tc.Close()

		tc.Get("/m/123").AssertRedirect(http.StatusSeeOther, "/setup")
	})

	t.Run("guests cannot view monitors", func(t *testing.T) {
		tc := NewTestCase(t)
		defer tc.Close()

		tc.CreateTestUser("test@example.com", "password")

		tc.Get("/m/123").AssertRedirect(http.StatusSeeOther, "/login")
	})

	t.Run("monitor has to exist", func(t *testing.T) {
		tc := NewTestCase(t)
		defer tc.Close()

		tc.LogIn().Get("/m/123").AssertStatusCode(http.StatusNotFound)
	})

	t.Run("monitor can be viewed", func(t *testing.T) {
		tc := NewTestCase(t)
		defer tc.Close()

		m, _ := tc.Store.CreateMonitor(t.Context(), uptimemonitor.Monitor{
			Url: "http://example.com",
		})

		tc.LogIn().Get(fmt.Sprintf("/m/%s", m.Uuid)).AssertStatusCode(http.StatusOK)
	})
}
