package testutil

import (
	"net/http"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

type AssertableResponse struct {
	T        *testing.T
	Response *http.Response
	Document *goquery.Document
}

func NewAssertableResponse(t *testing.T, res *http.Response) *AssertableResponse {
	doc, _ := goquery.NewDocumentFromReader(res.Body)

	return &AssertableResponse{
		T:        t,
		Response: res,
		Document: doc,
	}
}

func (ar *AssertableResponse) AssertStatusCode(expected int) *AssertableResponse {
	ar.T.Helper()

	if ar.Response.StatusCode != expected {
		ar.T.Fatalf("expected status code %d, got %d", expected, ar.Response.StatusCode)
	}

	return ar
}

func (ar *AssertableResponse) AssertRedirect(statusCode int, expected string) *AssertableResponse {
	ar.T.Helper()

	if ar.Response.Request.Response == nil {
		ar.T.Fatalf("no redirect has been made")
	}

	if ar.Response.Request.Response.StatusCode != statusCode {
		ar.T.Fatalf("expected status code %d, got %d", statusCode, ar.Response.Request.Response.StatusCode)
	}

	if ar.Response.Request.Response.Header.Get("Location") != expected {
		ar.T.Fatalf("expected redirect to %s, got %s", expected, ar.Response.Header.Get("Location"))
	}

	return ar
}

func (ar *AssertableResponse) AssertElementVisible(selector string) *AssertableResponse {
	ar.T.Helper()

	if ar.Document == nil {
		ar.T.Fatalf("no document available for assertion")
	}

	if ar.Document.Find(selector).Length() == 0 {
		html, _ := ar.Document.Html()
		ar.T.Fatalf(
			"expected element with selector '%s' to be visible but it was not found in the output: %v",
			selector,
			html,
		)
	}

	return ar
}

func (ar *AssertableResponse) AssertSeeText(text string) *AssertableResponse {
	ar.T.Helper()

	if ar.Document == nil {
		ar.T.Fatalf("no document available for assertion")
	}

	html, _ := ar.Document.Html()
	if !strings.Contains(html, text) {
		ar.T.Fatalf("expected to see text '%s' but it was not found in the output: %v", text, html)
	}

	return ar
}
