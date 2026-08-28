package store

import "testing"

func TestAddThenGet(t *testing.T) {
	s := New()
	code := s.Add("https://example.com")
	got, ok := s.Get(code)
	if !ok || got != "https://example.com" {
		t.Fatalf("Get(%q) = %q, %v", code, got, ok)
	}
}

func TestGetUnknownCode(t *testing.T) {
	if _, ok := New().Get("nope"); ok {
		t.Fatal("unknown code resolved")
	}
}
