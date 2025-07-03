package test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"uptimemonitor/handler"
	"uptimemonitor/pkg/testutil"
	"uptimemonitor/router"
	"uptimemonitor/store"
	"uptimemonitor/store/sqlite"
)

type TestCase struct {
	T      *testing.T
	Server *httptest.Server
	Client *http.Client
	Store  store.Store
}

func NewTestCase(t *testing.T) *TestCase {
	store := sqlite.New(":memory:")
	handler := handler.New(store)
	router := router.New(handler)
	server := httptest.NewServer(router)

	return &TestCase{
		T:      t,
		Server: server,
		Client: server.Client(),
		Store:  store,
	}
}

func (t *TestCase) Close() {
	t.Server.Close()
}

// todo ConfigurableRequest
func (t *TestCase) Get(url string) *testutil.AssertableResponse {
	res, err := t.Client.Get(t.Server.URL + url)
	if err != nil {
		t.T.Fatalf("failed to get %s: %v", url, err)
	}

	return testutil.NewAssertableResponse(t.T, res)
}
