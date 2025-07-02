package testutil

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"uptimemonitor/app"
	"uptimemonitor/database"
)

type TestCase struct {
	T      *testing.T
	App    *app.App
	Server *httptest.Server
	Client *http.Client
}

func NewTestCase(t *testing.T) *TestCase {
	db := database.Must(database.New(":memory:"))
	router := app.NewRouter(db)
	server := httptest.NewServer(router)

	return &TestCase{
		T:      t,
		App:    app.New(db),
		Server: server,
		Client: server.Client(),
	}
}

func (t *TestCase) Close() {
	t.Server.Close()
}

func (t *TestCase) Get(url string) *AssertableResponse {
	res, err := t.Client.Get(t.Server.URL + url)
	if err != nil {
		t.T.Fatalf("failed to get %s: %v", url, err)
	}

	return NewAssertableResponse(t.T, res)
}
