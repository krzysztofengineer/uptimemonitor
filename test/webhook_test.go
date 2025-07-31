package test

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"
	"time"
	"uptimemonitor"
	"uptimemonitor/form"
	"uptimemonitor/service"
)

func TestWebhook_SaveWebhook(t *testing.T) {
	t.Run("webhook is validated", func(t *testing.T) {
		tc := NewTestCase(t)
		defer tc.Close()

		f := form.MonitorForm{
			Url:        "https://valid.url",
			HttpMethod: http.MethodGet,

			HasWebhook:     true,
			WebhookUrl:     "invalid",
			WebhookHeaders: "invalid json",
			WebhookBody:    "data",
		}

		tc.AssertEqual(false, f.Validate())
	})

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

	t.Run("webhook data can be updated", func(t *testing.T) {
		tc := NewTestCase(t)
		defer tc.Close()

		_, err := tc.Store.CreateMonitor(t.Context(), uptimemonitor.Monitor{
			HttpMethod: http.MethodGet,
			Url:        "https://google.com",
		})

		tc.AssertNoError(err)

		tc.LogIn().
			Patch("/monitors/1", url.Values{
				"http_method":     []string{"GET"},
				"url":             []string{"https://example.com"},
				"has_webhook":     []string{"on"},
				"webhook_url":     []string{tc.Server.URL + "/test/webhook"},
				"webhook_method":  []string{"POST"},
				"webhook_headers": []string{`{"test":"abc"}`},
				"webhook_body":    []string{`{"test":"123"}`},
			}).AssertStatusCode(http.StatusOK)

		m, err := tc.Store.GetMonitorByID(t.Context(), 1)
		tc.AssertNoError(err)

		tc.AssertEqual(m.WebhookUrl, tc.Server.URL+"/test/webhook")
		tc.AssertEqual(m.WebhookMethod, "POST")
		tc.AssertEqual(m.WebhookHeaders, `{"test":"abc"}`)
		tc.AssertEqual(m.WebhookBody, `{"test":"123"}`)
	})

	t.Run("webhook fields are present in the forms", func(t *testing.T) {
		tc := NewTestCase(t)
		defer tc.Close()

		tc.LogIn().Get("/new").
			AssertElementVisible(`input[name="has_webhook"]`).
			AssertElementVisible(`select[name="webhook_method"]`).
			AssertElementVisible(`textarea[name="webhook_body"]`).
			AssertElementVisible(`textarea[name="webhook_headers"]`)

		m, _ := tc.Store.CreateMonitor(t.Context(), uptimemonitor.Monitor{
			HttpMethod: http.MethodGet,
			Url:        "https://google.com",
		})

		tc.Get(fmt.Sprintf("/m/%s/edit", m.Uuid)).
			AssertElementVisible(`input[name="has_webhook"]`).
			AssertElementVisible(`select[name="webhook_method"]`).
			AssertElementVisible(`textarea[name="webhook_body"]`).
			AssertElementVisible(`textarea[name="webhook_headers"]`)
	})

	t.Run("webhook is called on incident", func(t *testing.T) {
		tc := NewTestCase(t)
		defer tc.Close()

		ExpectedWebhookBody = `{"test":123}`
		ExpectedWebhookHeaderKey = "test"
		ExpectedWebhookHeaderValue = "abc"

		service := service.New(tc.Store)

		tc.Store.CreateMonitor(t.Context(), uptimemonitor.Monitor{
			Url:            tc.Server.URL + "/test/500",
			WebhookMethod:  http.MethodPost,
			WebhookUrl:     tc.Server.URL + "/test/webhook",
			WebhookHeaders: `{"test":"abc"}`,
			WebhookBody:    `{"test":123}`,
		})

		ch := service.StartCheck()
		service.RunCheck(t.Context(), ch)

		time.Sleep(time.Second * 3)

		tc.AssertDatabaseCount("incidents", 1)
		tc.AssertDatabaseCount("checks", 1)

		tc.AssertEqual(int64(TestWebhookCalledCount), int64(1))
	})

	t.Run("webhook can have parsed body", func(t *testing.T) {
		tc := NewTestCase(t)
		defer tc.Close()

		service := service.New(tc.Store)
		url := tc.Server.URL + "/test/500"

		tc.Store.CreateMonitor(t.Context(), uptimemonitor.Monitor{
			Url:           url,
			WebhookMethod: http.MethodPost,
			WebhookUrl:    tc.Server.URL + "/test/webhook",
			WebhookBody:   `{{.Url}},{{ .StatusCode}}`,
		})

		ch := service.StartCheck()
		service.RunCheck(t.Context(), ch)

		time.Sleep(time.Second * 3)

		tc.AssertDatabaseCount("incidents", 1)
		tc.AssertDatabaseCount("checks", 1)

		tc.AssertEqual(int64(TestWebhookCalledCount), int64(1))
		tc.AssertEqual(TestWebhookBody, fmt.Sprintf("%s,%d", url, 500))
	})
}
