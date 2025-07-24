package test

import (
	"net/http"
	"net/url"
	"testing"
)

func TestWebhook_SaveWebhook(t *testing.T) {
	t.Run("webhook info can be saved when creating a monitor", func(t *testing.T) {
		tc := NewTestCase(t)
		defer tc.Close()

		tc.LogIn().
			Post("/monitors", url.Values{
				"http_method":     []string{"GET"},
				"url":             []string{"https://example.com"},
				"has_webhook":     []string{"on"},
				"webhook_url":     []string{tc.Server.URL + "/test/webhook"},
				"webhook_method":  []string{"POST"},
				"webhook_headers": []string{`{"test":"abc"}`},
				"webhook_body":    []string{`{"test":"123"}`},
			}).AssertStatusCode(http.StatusOK)

		m, _ := tc.Store.GetMonitorByID(t.Context(), 1)

		tc.AssertEqual(m.WebhookUrl, tc.Server.URL+"/test/webhook")
		tc.AssertEqual(m.WebhookMethod, "POST")
		tc.AssertEqual(m.WebhookHeaders, `{"test":"abc"}`)
		tc.AssertEqual(m.WebhookBody, `{"test":"123"}`)
	})
}
