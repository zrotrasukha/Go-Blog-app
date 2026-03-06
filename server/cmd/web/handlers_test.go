package main

import (
	"net/http"
	"testing"

	"github.com/zrotrasukha/Go-Blog-app/internal/assert"
)

func TestHealth(t *testing.T) {
	app := newTestApplication(t)

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	statusCode, _, body := ts.get(t, "/health")

	assert.Equal(t, statusCode, http.StatusOK)
	assert.Equal(t, body, "OK")
}
