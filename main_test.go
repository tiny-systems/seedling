package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/tiny-systems/seedling/internal/store"
)

func TestShortenTrimsWhitespace(t *testing.T) {
	s := store.New()
	mux := newMux(s)

	form := url.Values{"url": {"  https://example.com  "}}
	req := httptest.NewRequest(http.MethodPost, "/shorten", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", rec.Code, rec.Body.String())
	}
	body, err := io.ReadAll(rec.Body)
	if err != nil {
		t.Fatal(err)
	}

	// Response is "http://<host>/<code>\n"; the code is the last path segment.
	response := strings.TrimSpace(string(body))
	code := response[strings.LastIndex(response, "/")+1:]

	got, ok := s.Get(code)
	if !ok || got != "https://example.com" {
		t.Fatalf("stored url = %q, %v; want %q", got, ok, "https://example.com")
	}
}
