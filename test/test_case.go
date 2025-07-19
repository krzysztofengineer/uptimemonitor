package test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
	"uptimemonitor"
	"uptimemonitor/handler"
	"uptimemonitor/pkg/testutil"
	"uptimemonitor/router"
	"uptimemonitor/service"
	"uptimemonitor/store"
	"uptimemonitor/store/sqlite"

	"golang.org/x/crypto/bcrypt"
)

type TestCase struct {
	T       *testing.T
	Server  *httptest.Server
	Client  *http.Client
	Store   store.Store
	User    *uptimemonitor.User
	Headers map[string]string
	Cookies []*http.Cookie
}

func NewTestCase(t *testing.T) *TestCase {
	store := sqlite.New(":memory:")
	service := service.New(store)
	handler := handler.New(store, service)
	router := router.New(handler)
	server := httptest.NewServer(router)

	return &TestCase{
		T:       t,
		Server:  server,
		Client:  server.Client(),
		Store:   store,
		Headers: map[string]string{},
		Cookies: []*http.Cookie{},
	}
}

func (tc *TestCase) Close() {
	tc.Server.Close()
}

func (tc *TestCase) WithHeader(key, value string) *TestCase {
	tc.Headers[key] = value

	return tc
}

func (tc *TestCase) WithCookie(c *http.Cookie) *TestCase {
	tc.Cookies = append(tc.Cookies, c)

	return tc
}

func (tc *TestCase) Get(url string) *testutil.AssertableResponse {
	req, err := http.NewRequest(http.MethodGet, tc.Server.URL+url, nil)
	if err != nil {
		tc.T.Fatalf("unexpected error: %v", err)
	}

	if len(tc.Cookies) > 0 {
		for _, c := range tc.Cookies {
			req.AddCookie(c)
		}
	}

	res, err := tc.Client.Do(req)
	if err != nil {
		tc.T.Fatalf("failed to get %s: %v", url, err)
	}

	return testutil.NewAssertableResponse(tc.T, res)
}

func (tc *TestCase) Post(url string, data url.Values) *testutil.AssertableResponse {
	req, err := http.NewRequest(http.MethodPost, tc.Server.URL+url, strings.NewReader(data.Encode()))
	if err != nil {
		tc.T.Fatalf("unexpected error: %v", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	if len(tc.Cookies) > 0 {
		for _, c := range tc.Cookies {
			req.AddCookie(c)
		}
	}

	res, err := tc.Client.Do(req)
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

func (tc *TestCase) LogIn() *TestCase {
	tc.CreateTestUser("test@example.com", "password")

	session, err := tc.Store.CreateSession(tc.T.Context(), uptimemonitor.Session{
		User:      *tc.User,
		UserID:    tc.User.ID,
		ExpiresAt: time.Now().Add(time.Hour),
	})
	if err != nil {
		tc.T.Fatalf("unexpected error: %v", err)
	}

	c := &http.Cookie{
		Name:     "session",
		Value:    session.Uuid,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   false,
		Expires:  session.ExpiresAt,
	}

	return tc.WithCookie(c)
}
