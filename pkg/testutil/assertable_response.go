package testutil

import (
	"net/http"
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

func (r *AssertableResponse) AssertStatusCode(expected int) *AssertableResponse {
	r.T.Helper()

	if r.Response.StatusCode != expected {
		r.T.Fatalf("expected status code %d, got %d", expected, r.Response.StatusCode)
	}

	return r
}

func (r *AssertableResponse) AssertRedirect(statusCode int, expected string) *AssertableResponse {
	r.T.Helper()

	if r.Response.Request.Response == nil {
		r.T.Fatalf("no redirect has been made")
	}

	if r.Response.Request.Response.StatusCode != statusCode {
		r.T.Fatalf("expected status code %d, got %d", statusCode, r.Response.Request.Response.StatusCode)
	}

	if r.Response.Request.Response.Header.Get("Location") != expected {
		r.T.Fatalf("expected redirect to %s, got %s", expected, r.Response.Header.Get("Location"))
	}

	return r
}

func (r *AssertableResponse) AssertElementVisible(selector string) *AssertableResponse {
	r.T.Helper()

	if r.Document == nil {
		r.T.Fatalf("no document available for assertion")
	}

	if r.Document.Find(selector).Length() == 0 {
		html, _ := r.Document.Html()
		r.T.Fatalf(
			"expected element with selector '%s' to be visible but it was not found in the output: %v",
			selector,
			html,
		)
	}

	return r
}
