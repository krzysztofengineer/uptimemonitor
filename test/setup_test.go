package test

import (
	"net/http"
	"testing"
)

func TestSetup(t *testing.T) {
	t.Run("redirects to setup when no users are found", func(t *testing.T) {
		test := NewTestCase(t)
		defer test.Close()

		test.Get("/").
			AssertRedirect(http.StatusSeeOther, "/setup")
	})
}
