package test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"uptimemonitor"
	"uptimemonitor/handler"
	"uptimemonitor/pkg/testutil"
	"uptimemonitor/router"
	"uptimemonitor/store"
	"uptimemonitor/store/sqlite"

	"golang.org/x/crypto/bcrypt"
)

type TestCase struct {
	T      *testing.T
	Server *httptest.Server
	Client *http.Client
	Store  store.Store
	User   *uptimemonitor.User
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

func (tc *TestCase) Close() {
	tc.Server.Close()
}

func (tc *TestCase) Get(url string) *testutil.AssertableResponse {
	res, err := tc.Client.Get(tc.Server.URL + url)
	if err != nil {
		tc.T.Fatalf("failed to get %s: %v", url, err)
	}

	return testutil.NewAssertableResponse(tc.T, res)
}

func (tc *TestCase) Post(url string, data url.Values) *testutil.AssertableResponse {
	res, err := tc.Client.Post(tc.Server.URL+url, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
	if err != nil {
		tc.T.Fatalf("failed to post %s: %v", url, err)
	}

	return testutil.NewAssertableResponse(tc.T, res)
}

func (tc *TestCase) AssertDatabaseCount(table string, expected int) *TestCase {
	tc.T.Helper()

	stmt := fmt.Sprintf(`SELECT COUNT(*) FROM %s`, table)
	var count int

	err := tc.Store.DB().QueryRow(stmt).Scan(&count)
	if err != nil {
		tc.T.Fatalf("failed to count rows from table '%s', error: %v", table, err)
	}

	if count != expected {
		tc.T.Fatalf("expected to find %d number of rows in a table '%s, but found %d", expected, table, count)
	}

	return tc
}

func (tc *TestCase) CreateTestUser(email, password string) *TestCase {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		tc.T.Fatalf("unexpected bcrypt error: %v", err)
	}

	user, err := tc.Store.CreateUser(tc.T.Context(), uptimemonitor.User{
		Name:         "Test User",
		Email:        email,
		PasswordHash: string(hash),
	})
	if err != nil {
		tc.T.Fatalf("unable to create test user: %v", err)
	}

	tc.User = &user

	return tc
}
