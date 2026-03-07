package main

import (
	"bytes"
	"html"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"regexp"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form"
	"github.com/zrotrasukha/Go-Blog-app/internal/models/mocks"
)

func (ts *testServer) postForm(t *testing.T, urlPath string, formData url.Values) (int, http.Header, string) {
	rs, err := ts.Client().PostForm(ts.URL+urlPath, formData)
	if err != nil {
		t.Fatal(err)
	}

	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	bytes.TrimSpace(body)
	if err != nil {
		t.Fatal(err)
	}

	return rs.StatusCode, rs.Header, string(body)
}

var csrfRGXExp = regexp.MustCompile(`name=['"]csrf_token['"][^>]*value=['"]([^'"]+)['"]`)

func extractCSRFToken(t *testing.T, body string) string {
	matches := csrfRGXExp.FindSubmatch([]byte(body))
	if len(matches) < 2 {
		t.Fatal("no csrf token found in body")
	}

	return html.UnescapeString(string(matches[1]))
}

func newTestApplication(t *testing.T) *application {
	templateCache, err := newTemplateCache()
	if err != nil {
		t.Fatal(err)
	}

	formDecoder := form.NewDecoder()
	sessionManager := scs.New()
	sessionManager.Lifetime = 12 * time.Hour

	return &application{
		errorlog:       log.New(io.Discard, "", 0),
		infolog:        log.New(io.Discard, "", 0),
		blogs:          &mocks.BlogModel{},
		users:          &mocks.UserModel{},
		sessionManager: sessionManager,
		formDecoder:    formDecoder,
		templateCache:  templateCache,
	}
}

type testServer struct {
	*httptest.Server
}

func newTestServer(t *testing.T, h http.Handler) *testServer {
	ts := httptest.NewTLSServer(h)

	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}

	ts.Client().Jar = jar

	ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	return &testServer{ts}
}

// get is a helper method which makes a GET request to the test server and returns the response status code, headers and body.
func (ts *testServer) get(t *testing.T, urlPath string) (int, http.Header, string) {
	rs, err := ts.Client().Get(ts.URL + urlPath)
	if err != nil {
		t.Fatal(err)
	}

	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}

	bytes.TrimSpace(body)

	return rs.StatusCode, rs.Header, string(body)

}
