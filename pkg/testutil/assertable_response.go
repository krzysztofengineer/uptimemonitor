package testutil

import (
	"net/http"
	"testing"
)

type AssertableResponse struct {
	T        *testing.T
	Response *http.Response
}

func NewAssertableResponse(t *testing.T, res *http.Response) *AssertableResponse {
	return &AssertableResponse{
		T:        t,
		Response: res,
	}
}

func (r *AssertableResponse) AssertStatusCode(expected int) {
	r.T.Helper()

	if r.Response.StatusCode != expected {
		r.T.Fatalf("expected status code %d, got %d", expected, r.Response.StatusCode)
	}
}

func (r *AssertableResponse) AssertRedirect(statusCode int, expected string) {
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
}
