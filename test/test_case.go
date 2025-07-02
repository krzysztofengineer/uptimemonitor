package test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"uptimemonitor/app"
	"uptimemonitor/pkg/testutil"
	"uptimemonitor/sqlite"
)

type TestCase struct {
	T      *testing.T
	Server *httptest.Server
	Client *http.Client
}

func NewTestCase(t *testing.T) *TestCase {
	store := sqlite.New(":memory:")
	handler := app.NewHandler(store)
	router := app.NewRouter(handler)
	server := httptest.NewServer(router)

	return &TestCase{
		T:      t,
		Server: server,
		Client: server.Client(),
	}
}

func (t *TestCase) Close() {
	t.Server.Close()
}

func (t *TestCase) Get(url string) *testutil.AssertableResponse {
	res, err := t.Client.Get(t.Server.URL + url)
	if err != nil {
		t.T.Fatalf("failed to get %s: %v", url, err)
	}

	return testutil.NewAssertableResponse(t.T, res)
}
