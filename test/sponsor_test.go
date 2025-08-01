package test

import (
	"net/http"
	"testing"
)

func TestSponsor(t *testing.T) {
	t.Run("sponsors are lazy loaded", func(t *testing.T) {
		tc := NewTestCase(t)
		defer tc.Close()

		tc.Get("/sponsors").
			AssertStatusCode(http.StatusOK).
			AssertElementVisible(`div[hx-get="/sponsors"]`)
	})

	t.Run("sponsors are loaded via api", func(t *testing.T) {
		tc := NewTestCase(t)
		defer tc.Close()

		tc.WithHeader("HX-Request", "true").
			Get("/sponsors").
			AssertSeeText("AIR Labs")
	})
}
