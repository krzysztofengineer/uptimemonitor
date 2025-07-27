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
}
