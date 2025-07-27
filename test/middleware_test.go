package test

import (
	"net/http"
	"testing"
)

func TestMiddleware(t *testing.T) {
	t.Run("panic recoverer", func(t *testing.T) {
		tc := NewTestCase(t)
		defer tc.Close()

		tc.Get("/test/panic").AssertStatusCode(http.StatusInternalServerError)
	})

	t.Run("cache test", func(t *testing.T) {
		tc := NewTestCase(t)
		defer tc.Close()

		tc.LogIn().
			Get("/").
			AssertStatusCode(http.StatusOK).
			AssertHeader("Cache-Control", "no-cache, no-store, must-revalidate").
			AssertHeader("Pragma", "no-cache").
			AssertHeader("Expires", "0")
	})
}
