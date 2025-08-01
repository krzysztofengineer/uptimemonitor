package test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"
	"time"
	"uptimemonitor"
	"uptimemonitor/handler"
	"uptimemonitor/pkg/testutil"
	"uptimemonitor/router"
	"uptimemonitor/service"
	"uptimemonitor/store"

	"golang.org/x/crypto/bcrypt"
)

var TestWebhookCalledCount int64
var TestWebhookBody string
var ExpectedWebhookBody string
var ExpectedWebhookHeaderKey string
var ExpectedWebhookHeaderValue string

type TestCase struct {
	T       *testing.T
	Server  *httptest.Server
	Client  *http.Client
	Store   *store.Store
	User    *uptimemonitor.User
	Headers map[string]string
	Cookies []*http.Cookie
}

func NewTestCase(t *testing.T) *TestCase {
	store := store.New(":memory:")
	service := service.New(store)
	handler := handler.New(store, service, false)
	router := router.New(handler, registerRoutes)
	server := httptest.NewServer(router)

	TestWebhookBody = ""
	TestWebhookCalledCount = 0
	ExpectedWebhookBody = ""
	ExpectedWebhookHeaderKey = ""
	ExpectedWebhookHeaderValue = ""

	return &TestCase{
		T:       t,
		Server:  server,
		Client:  server.Client(),
		Store:   store,
		Headers: map[string]string{},
		Cookies: []*http.Cookie{},
	}
}

func registerRoutes(router *http.ServeMux) {
	router.HandleFunc("GET /test/200", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	router.HandleFunc("GET /test/404", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	router.HandleFunc("GET /test/500", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})

	router.HandleFunc("GET /test/timeout", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(30 * time.Second)
		w.WriteHeader(http.StatusOK)
	})

	router.HandleFunc("GET /test/panic", func(w http.ResponseWriter, r *http.Request) {
		panic("test")
	})

	i := 0
	router.HandleFunc("GET /test/even", func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			i++
		}()

		if i%2 == 0 {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	router.HandleFunc("POST /test/post", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	router.HandleFunc("PATCH /test/patch", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	router.HandleFunc("PUT /test/put", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	router.HandleFunc("DELETE /test/delete", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	router.HandleFunc("POST /test/body", func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		if string(body) != `{"test":123}` {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	})

	router.HandleFunc("POST /test/headers", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("test") != "abc" {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	})

	router.HandleFunc("POST /test/webhook", func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		if ExpectedWebhookBody != "" && string(body) != ExpectedWebhookBody {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		if ExpectedWebhookHeaderKey != "" && r.Header.Get(ExpectedWebhookHeaderKey) != ExpectedWebhookHeaderValue {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		TestWebhookCalledCount++
		TestWebhookBody = string(body)

		w.WriteHeader(http.StatusOK)
	})
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

	if len(tc.Headers) > 0 {
		for k, v := range tc.Headers {
			req.Header.Set(k, v)
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

	for _, c := range res.Cookies() {
		if c.Name == "session" {
			tc.Cookies = append(tc.Cookies, c)
		}
	}

	return testutil.NewAssertableResponse(tc.T, res)
}

func (tc *TestCase) Patch(url string, data url.Values) *testutil.AssertableResponse {
	req, err := http.NewRequest(http.MethodPatch, tc.Server.URL+url, strings.NewReader(data.Encode()))
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

func (tc *TestCase) Delete(url string) *testutil.AssertableResponse {
	req, err := http.NewRequest(http.MethodDelete, tc.Server.URL+url, nil)
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
		tc.T.Fatalf("failed to delete %s: %v", url, err)
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

func (tc *TestCase) AssertEqual(a, b any) *TestCase {
	tc.T.Helper()

	if !reflect.DeepEqual(a, b) {
		tc.T.Fatalf(`expected "%v" to be equal to "%v"`, a, b)
	}
	return tc
}

func (tc *TestCase) AssertNoError(err error) *TestCase {
	tc.T.Helper()

	if err != nil {
		tc.T.Fatalf("unexpected error: %v", err)
	}

	return tc
}
